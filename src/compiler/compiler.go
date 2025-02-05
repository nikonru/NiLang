package compiler

import (
	"NiLang/src/ast"
	"NiLang/src/helper"
	"NiLang/src/lexer"
	"NiLang/src/parser"
	"NiLang/src/tokens"
	"bytes"
	"fmt"
	"log"
	"strconv"
)

type errors = []helper.Error

type Compiler struct {
	output           bytes.Buffer
	memoryIndex      address
	stackMemoryIndex address

	scope *scope

	labelIndex uint64

	maxStackAddress address
	errors          errors
}

func New(stackSize int) *Compiler {

	return &Compiler{
		memoryIndex:      address(stackSize),
		stackMemoryIndex: -1,
		scope:            newScope(""),
		labelIndex:       0,
		maxStackAddress:  address(stackSize)}
}

func (c *Compiler) Compile(input []byte) ([]byte, errors) {

	lexer := lexer.New(input)
	parser := parser.New(&lexer)
	program := parser.Parse()
	if program == nil {
		log.Fatalf("parser.Parse() has returned nil")
	}

	errors := parser.Errors()
	if len(errors) != 0 {
		fmt.Printf("parser had %d error(s)\n", len(errors))
		fmt.Print("parser error(s):\n")

		return c.output.Bytes(), errors
	}

	fmt.Println("PROGRAM TREE")
	for _, statement := range program.Statements {
		fmt.Println(statement.String())
		c.compileStatement(statement)
	}
	fmt.Println("END")
	fmt.Println(c.scope.variables)
	return c.output.Bytes(), c.errors
}

func (c *Compiler) emit(op command, args ...interface{}) {
	write := func(arg interface{}, id int) {
		switch v := arg.(type) {
		case int:
			c.output.WriteString(strconv.Itoa(v))
		case int64:
			c.output.WriteString(strconv.FormatInt(v, 10))
		case string:
			c.output.WriteString(v)
		case bool:
			if v {
				c.output.WriteString("1")
			} else {
				c.output.WriteString("0")
			}
		default:
			log.Fatalf("type of arg%d not handled. got=%T", id, arg)
		}
	}

	c.output.WriteString(op)
	for i, arg := range args {
		c.output.WriteString(" ")
		write(arg, i)
	}

	c.output.WriteString("\n")
}

func (c *Compiler) emitLabel(label string) {
	c.emit(label + ":")
}

func (c *Compiler) compileStatement(statement ast.Statement) {
	switch stm := statement.(type) {
	case *ast.DeclarationStatement:
		c.compileDeclarationStatement(stm)
	case *ast.ExpressionStatement:
		c.compileExpression(stm.Expression)
	case *ast.ReturnStatement:
		c.compileReturnStatement(stm)
	case *ast.UsingStatement:
		c.compileUsingStatement(stm)
	case *ast.AssignmentStatement:
		c.compileAssignmentStatement(stm)
	case *ast.ScopeStatement:
		c.compileScopeStatement(stm)
	default:
		log.Fatalf("type of statement is not handled. got=%T", statement)
	}
	c.flushStackMemory()
}

func (c *Compiler) compileDeclarationStatement(ds *ast.DeclarationStatement) {
	_type, register := c.compileExpression(ds.Value)

	if _type != ds.Var.Type {
		err := helper.MakeError(ds.Var.Token, fmt.Sprintf("declared variable and expression have different types. variable=%q, expression=%q", ds.Var.Type, _type))
		c.addError(err)
	}

	addr := c.purchaseMemoryAddress()
	c.emit(LOAD_TO_MEM_FROM_REG, addr, register)

	if ok := c.scope.AddVariable(ds.Var.Name, addr, _type); !ok {
		err := helper.MakeError(ds.Var.Token, fmt.Sprintf("redeclaration of variable %q", ds.Var.Name))
		c.addError(err)
	}
}

func (c *Compiler) compileReturnStatement(rs *ast.ReturnStatement) {
	returnType, ok := c.scope.returnType.(string)

	if !ok {
		err := helper.MakeError(rs.Token, "unexpected return statement")
		c.addError(err)
	}

	var _type, register name
	if rs.Value != nil {
		_type, register = c.compileExpression(rs.Value)
	}

	if returnType != _type {
		err := helper.MakeError(rs.Token, fmt.Sprintf("expected return of type=%q, got=%q", returnType, _type))
		c.addError(err)
	}

	c.emit(LOAD_TO_REG_FROM_REG, AX, register) // TODO: maybe we can select some area of memory for this
	c.emit(RETURN)
}

func (c *Compiler) compileUsingStatement(us *ast.UsingStatement) {
	switch name := us.Name.(type) {
	case *ast.Identifier:
		c.scope.UsingScope(name.Value)
	case *ast.ScopeExpression:
		log.Fatal("WIP") //TODO implement scope expression compiling and use it here
	default:
		err := helper.MakeError(us.Token, fmt.Sprintf("expected identifier or scope expression of scope, got=%T", name))
		c.addError(err)
	}
}

func (c *Compiler) compileAssignmentStatement(as *ast.AssignmentStatement) {
	_type, register := c.compileExpression(as.Value)
	variable, ok := c.scope.GetVariable(as.Name.Value)

	if !ok {
		err := helper.MakeError(as.Name.Token, fmt.Sprintf("assigning to undeclared variable '%q'", as.Name.Value))
		c.addError(err)
	}

	if variable.Type != _type {
		err := helper.MakeError(as.Name.Token, fmt.Sprintf("expected expression of type=%q, got=%q", variable.Type, _type))
		c.addError(err)
	}

	c.emit(LOAD_TO_MEM_FROM_REG, variable.Addr, register)
}

func (c *Compiler) compileScopeStatement(ss *ast.ScopeStatement) {

	c.enterNamedScope(ss.Name.Value)

	if ok := c.scope.GetParent().AddScope(c.scope); !ok {
		err := helper.MakeError(ss.Name.Token, fmt.Sprintf("redeclaration of scope %q", c.scope.name))
		c.addError(err)
	}

	for _, statement := range ss.Body.Statements {
		c.compileStatement(statement)
	}

	c.leaveScope()
}

func (c *Compiler) compileExpression(statement ast.Expression) (name, register) {
	switch exp := statement.(type) {
	case *ast.IntegralLiteral:
		return c.compileIntegralLiteral(exp)
	case *ast.BooleanLiteral:
		return c.compileBooleanLiteral(exp)
	case *ast.PrefixExpression:
		return c.compilePrefixExpression(exp)
	case *ast.InfixExpression:
		return c.compileInfixExpression(exp)
	case *ast.Identifier:
		return c.compileIdentifier(exp)
	case *ast.CallExpression:
		log.Fatalf("WIP")
		return "", ""
	case *ast.ScopeExpression:
		log.Fatalf("WIP")
		return "", ""
	default:
		log.Fatalf("type of expression is not handled. got=%T", exp)
		return "", ""
	}
}

func (c *Compiler) compileIntegralLiteral(expression *ast.IntegralLiteral) (name, register) {
	c.emit(LOAD_TO_REG_FROM_VAL, AX, expression.Value)
	return Int, AX
}

func (c *Compiler) compileBooleanLiteral(expression *ast.BooleanLiteral) (name, register) {
	value := BOOL_FALSE
	if expression.Value {
		value = BOOL_TRUE
	}

	c.emit(LOAD_TO_REG_FROM_VAL, AX, value)
	return Bool, AX
}

func (c *Compiler) compilePrefixExpression(expression *ast.PrefixExpression) (name, register) {

	_type, register := c.compileExpression(expression.Right)

	switch expression.Operator {
	case tokens.NOT:
		if _type != Bool {
			err := helper.MakeError(expression.Token, fmt.Sprintf("expected boolean expression. got=%q", _type))
			c.addError(err)
		}

		end := c.getUniqueLabel()
		True := c.getUniqueLabel()
		False := c.getUniqueLabel()

		c.emit(COMPARE_WITH_VALUE, register, BOOL_TRUE)
		c.emit(JUMP_IF_EQUAL, register, False)

		c.emitLabel(True)
		c.emit(LOAD_TO_REG_FROM_VAL, register, BOOL_TRUE)
		c.emit(JUMP, end)

		c.emitLabel(False)
		c.emit(LOAD_TO_REG_FROM_VAL, register, BOOL_FALSE)

		c.emitLabel(end)

	default:
		log.Fatalf("type of prefix is not handled. got=%q", expression.Operator)
	}

	return Bool, register
}

func (c *Compiler) compileInfixExpression(expression *ast.InfixExpression) (name, register) {

	leftType, leftRegister := c.compileExpression(expression.Left)
	buffer := c.purchaseStackMemoryAddress()
	c.emit(LOAD_TO_MEM_FROM_REG, buffer, leftRegister)

	rightType, rightRegister := c.compileExpression(expression.Right)

	if rightRegister != BX {
		c.emit(LOAD_TO_REG_FROM_REG, BX, rightRegister)
		rightRegister = BX
	}

	c.emit(LOAD_TO_REG_FROM_MEM, AX, buffer)
	leftRegister = AX

	emitComparison := func(jump command) (name, register) {
		if leftType != Int || rightType != Int {
			err := helper.MakeError(expression.Token, fmt.Sprintf("expected int expression(s). got left=%q and right=%q", leftType, rightType))
			c.addError(err)
		}

		end := c.getUniqueLabel()
		True := c.getUniqueLabel()

		c.emit(COMPARE, leftRegister, rightRegister)
		c.emit(jump, True)
		c.emit(LOAD_TO_REG_FROM_VAL, AX, BOOL_FALSE)
		c.emit(JUMP, end)

		c.emitLabel(True)
		c.emit(LOAD_TO_REG_FROM_VAL, AX, BOOL_TRUE)

		c.emitLabel(end)
		return Bool, AX
	}

	switch expression.Operator {
	case tokens.LT:
		return emitComparison(JUMP_IF_LESS_THAN)
	case tokens.LE:
		return emitComparison(JUMP_IF_LESS_EQUAL_THAN)
	case tokens.GT:
		return emitComparison(JUMP_IF_GREATER_THAN)
	case tokens.GE:
		return emitComparison(JUMP_IF_GREATER_EQUAL_THAN)
	case tokens.NEQUAL:
		return emitComparison(JUMP_IF_NOT_EQUAL)
	case tokens.EQUAL:
		return emitComparison(JUMP_IF_EQUAL)
	case tokens.AND:
		if leftType != Bool || rightType != Bool {
			err := helper.MakeError(expression.Token, fmt.Sprintf("expected bool expression(s). got left=%q and right=%q", leftType, rightType))
			c.addError(err)
		}

		end := c.getUniqueLabel()
		False := c.getUniqueLabel()

		c.emit(COMPARE_WITH_VALUE, leftRegister, BOOL_FALSE)
		c.emit(JUMP_IF_EQUAL, False)
		c.emit(COMPARE_WITH_VALUE, rightRegister, BOOL_FALSE)
		c.emit(JUMP_IF_EQUAL, False)

		c.emit(LOAD_TO_REG_FROM_VAL, AX, BOOL_TRUE)
		c.emit(JUMP, end)

		c.emitLabel(False)
		c.emit(LOAD_TO_REG_FROM_VAL, AX, BOOL_FALSE)

		c.emitLabel(end)
		return Bool, AX
	case tokens.OR:
		if leftType != Bool || rightType != Bool {
			err := helper.MakeError(expression.Token, fmt.Sprintf("expected bool expression(s). got left=%q and right=%q", leftType, rightType))
			c.addError(err)
		}

		end := c.getUniqueLabel()
		True := c.getUniqueLabel()

		c.emit(COMPARE_WITH_VALUE, leftRegister, BOOL_TRUE)
		c.emit(JUMP_IF_EQUAL, True)
		c.emit(COMPARE_WITH_VALUE, rightRegister, BOOL_TRUE)
		c.emit(JUMP_IF_EQUAL, True)

		c.emit(LOAD_TO_REG_FROM_VAL, AX, BOOL_FALSE)
		c.emit(JUMP, end)

		c.emitLabel(True)
		c.emit(LOAD_TO_REG_FROM_VAL, AX, BOOL_TRUE)

		c.emitLabel(end)
		return Bool, AX
	default:
		log.Fatalf("type of infix expression is not handled. got=%q", expression.Operator)
		return "", ""
	}
}

func (c *Compiler) compileIdentifier(expression *ast.Identifier) (name, register) {
	if variable, ok := c.scope.GetVariable(expression.Value); ok {
		c.emit(LOAD_TO_REG_FROM_MEM, AX, variable.Addr)
		return variable.Type, AX
	}

	err := helper.MakeError(expression.Token, fmt.Sprintf("unknown identifier. got=%q", expression))
	c.addError(err)
	return "", ""
}

func (c *Compiler) purchaseMemoryAddress() address {
	c.memoryIndex++
	return c.memoryIndex
}

func (c *Compiler) purchaseStackMemoryAddress() address {
	c.stackMemoryIndex++
	if c.stackMemoryIndex >= c.maxStackAddress {
		log.Fatalf("Stack overflow, StackSize=%d", c.maxStackAddress)
	}
	return c.stackMemoryIndex
}

func (c *Compiler) flushStackMemory() {
	c.stackMemoryIndex = -1
}

func (c *Compiler) getUniqueLabel() string {
	// TODO: maximize the number of possible labels
	c.labelIndex++
	return "label" + strconv.FormatUint(c.labelIndex, 10)
}

func (c *Compiler) addError(error helper.Error) {
	c.errors = append(c.errors, error)
}

func (c *Compiler) enterNamedScope(name name) {
	scope := newScope(name)
	scope.SetParent(c.scope)
	c.scope = scope
}

func (c *Compiler) enterScope() {
	c.enterNamedScope("")
}

func (c *Compiler) leaveScope() {
	parent := c.scope.GetParent()
	if parent != nil {
		c.scope = parent
	} else {
		log.Fatal("leaving global scope")
	}
}
