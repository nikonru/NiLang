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
	"slices"
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
	fmt.Println(c.scope.functions)
	fmt.Println(c.scope.children)
	for _, scope := range c.scope.usingScopes {
		fmt.Println(scope.name)
	}

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
	case *ast.WhileStatement:
		c.compileWhileStatement(stm)
	case *ast.AliasStatement:
		c.compileAliasStatement(stm)
	case *ast.FunctionStatement:
		c.compileFunctionStatement(stm)
	case *ast.IfStatement:
		c.compileIfStatement(stm)
	case *ast.BreakStatement:
		c.compileBreakStatement(stm)
	case *ast.ContinueStatement:
		c.compileContinueStatement(stm)
	default:
		log.Fatalf("type of statement is not handled. got=%T", statement)
	}
	c.flushStackMemory()
}

func (c *Compiler) compileDeclarationStatement(ds *ast.DeclarationStatement) {
	_type, register := c.compileExpression(ds.Value)

	var_type, ok := c.findType(&ds.Var)
	if !ok {
		return
	}

	if _type != var_type {
		err := helper.MakeError(ds.Var.Token, fmt.Sprintf("declared variable and expression have different types. variable=%q, expression=%q",
			ds.Var.Type, _type.String()))
		c.addError(err)
	}

	if ok := c.addNewVariable(register, ds.Var.Name, var_type); !ok {
		err := helper.MakeError(ds.Var.Token, fmt.Sprintf("redeclaration of variable %q", ds.Var.Name))
		c.addError(err)
	}
}

func (c *Compiler) addNewVariable(register register, name name, t Type) bool {
	addr := c.purchaseMemoryAddress()
	c.emit(LOAD_TO_MEM_FROM_REG, addr, register)

	return c.scope.AddVariable(name, addr, t)
}

func (c *Compiler) compileReturnStatement(rs *ast.ReturnStatement) {
	returnType, ok := c.scope.returnType.(Type)

	if !ok {
		err := helper.MakeError(rs.Token, "unexpected return statement")
		c.addError(err)
	}

	var register name
	_type := VOID
	if rs.Value != nil {
		_type, register = c.compileExpression(rs.Value)
	}

	if returnType != _type {
		err := helper.MakeError(rs.Token, fmt.Sprintf("expected return of type=%q, got=%q",
			returnType.String(), _type.String()))
		c.addError(err)
	}

	if register != "" {
		c.emit(LOAD_TO_REG_FROM_REG, RETURN_REGISTER, register)
	}
	c.emit(RETURN)
}

func (c *Compiler) compileUsingStatement(us *ast.UsingStatement) {
	var s *scope
	var ok bool

	switch name := us.Name.(type) {
	case *ast.Identifier:
		s, ok = c.scope.GetScope(name.Value)
	case *ast.ScopeExpression:
		s, ok = c.findScope(name, c.scope)
		if !ok {
			err := helper.MakeError(us.Token, "undeclared scope/alias expression")
			c.addError(err)
		}
		s, ok = s.GetScope(name.Value.Value)
	default:
		err := helper.MakeError(us.Token, fmt.Sprintf("expected identifier or scope expression of scope, got=%T", name))
		c.addError(err)
	}

	if !ok {
		err := helper.MakeError(us.Token, "undeclared scope/alias expression")
		c.addError(err)
	} else {
		c.scope.UsingScope(s)
	}
}

func (c *Compiler) compileAssignmentStatement(as *ast.AssignmentStatement) {
	_type, register := c.compileExpression(as.Value)
	variable, ok := c.scope.GetVariable(as.Name.Value)

	if !ok {
		err := helper.MakeError(as.Name.Token, fmt.Sprintf("assigning to undeclared variable %q", as.Name.Value))
		c.addError(err)
	}

	if variable.Type != _type {
		err := helper.MakeError(as.Name.Token, fmt.Sprintf("expected expression of type=%q, got=%q",
			variable.Type.String(), _type.String()))
		c.addError(err)
	}

	c.emit(LOAD_TO_MEM_FROM_REG, variable.Addr, register)
}

func (c *Compiler) compileScopeStatement(ss *ast.ScopeStatement) {
	c.enterNamedScope(ss.Name.Value)
	defer c.leaveScope()

	if ok := c.scope.GetParent().AddScope(c.scope); !ok {
		err := helper.MakeError(ss.Name.Token, fmt.Sprintf("redeclaration of scope/alias %q", c.scope.name))
		c.addError(err)
	}

	for _, statement := range ss.Body.Statements {
		c.compileStatement(statement)
	}
}

func (c *Compiler) compileWhileStatement(ws *ast.WhileStatement) {

	loop := c.getUniqueLabel()
	end := c.getUniqueLabel()

	c.emitLabel(loop)
	_type, register := c.compileExpression(ws.Condition)

	if _type != builtIn(Bool) {
		err := helper.MakeError(ws.Token, fmt.Sprintf("expected boolean condition in while loop, got %q", _type.String()))
		c.addError(err)
	}

	c.emit(COMPARE_WITH_VALUE, register, BOOL_TRUE)
	c.emit(JUMP_IF_NOT_EQUAL, end)

	c.enterScope()
	defer c.leaveScope()

	c.scope.escapeLabel = end
	c.scope.repeatLabel = loop

	for _, statement := range ws.Body.Statements {
		c.compileStatement(statement)
	}

	c.emit(JUMP, loop)
	c.emitLabel(end)
}

func (c *Compiler) compileAliasStatement(as *ast.AliasStatement) {
	c.enterNamedScope(helper.FirstToLowerCase(as.Var.Name))
	defer c.leaveScope()

	if ok := c.scope.GetParent().AddScope(c.scope); !ok {
		err := helper.MakeError(as.Token, fmt.Sprintf("redeclaration of scope/alias %q", c.scope.name))
		c.addError(err)
	}

	t, ok := as.Var.Type.(*ast.Identifier)
	if !ok || (t.Value != Bool && t.Value != Int) {
		err := helper.MakeError(as.Token, fmt.Sprintf("expected alias to be primitive type(Bool, Int), got %q", as.Var.Type))
		c.addError(err)
	} else {
		for _, val := range as.Values {
			switch v := val.Value.(type) {
			case *ast.IntegralLiteral, *ast.BooleanLiteral:
				_type, register := c.compileExpression(val.Value)

				if _type.Name != t.Value {
					err := helper.MakeError(val.Var.Token, fmt.Sprintf("declared alias and expression have different types. alias=%q, expression=%q",
						as.Var.Type, _type.String()))
					c.addError(err)
				}

				if ok := c.addNewVariable(register, val.Var.Name, Type{Scope: c.scope.GetParent(), Name: as.Var.Name}); !ok {
					err := helper.MakeError(val.Var.Token, fmt.Sprintf("redeclaration of alias %q", val.Var.Name))
					c.addError(err)
				}
			default:
				err := helper.MakeError(val.Var.Token, fmt.Sprintf("expected literal expression, got %T", v))
				c.addError(err)
			}
		}
	}
}

func (c *Compiler) compileFunctionStatement(fs *ast.FunctionStatement) {
	_type := VOID

	if fs.Var.Type != nil {
		var ok bool
		_type, ok = c.findType(&fs.Var)

		if !ok {
			err := helper.MakeError(fs.Var.Token, "undeclared function type")
			c.addError(err)
		}
	}

	start := c.getUniqueLabel()
	end := c.getUniqueLabel()

	var arguments []variable
	if fs.Parameters != nil {
		arguments = make([]variable, len(fs.Parameters))

		for i, parameter := range fs.Parameters {
			if parameter.Type != nil {
				var _var variable
				_var.Name = parameter.Name
				_type, ok := c.findType(&parameter)
				if !ok {
					err := helper.MakeError(parameter.Token, "undeclared parameter type")
					c.addError(err)
					return
				}
				_var.Type = _type
				_var.Addr = c.purchaseMemoryAddress()

				arguments[i] = _var
			} else {
				err := helper.MakeError(parameter.Token, "undeclared parameter type")
				c.addError(err)
				return
			}
		}
	} else {
		arguments = make([]variable, 0)
	}

	ok := c.scope.AddFunction(fs.Var.Name, start, _type, arguments)
	if !ok {
		err := helper.MakeError(fs.Token, fmt.Sprintf("redeclaration of function %q", fs.Var.Name))
		c.addError(err)
		return
	}

	c.enterNamedScope(fs.Var.Name)
	defer c.leaveScope()
	c.scope.returnType = _type

	for i, arg := range arguments {
		ok = c.scope.AddVariable(arg.Name, arg.Addr, arg.Type)
		if !ok {
			err := helper.MakeError(fs.Parameters[i].Token, fmt.Sprintf("redeclaration of an argument %q", arg.Name))
			c.addError(err)
		}
	}

	c.emit(JUMP, end)
	c.emitLabel(start)

	foundReturnStatement := false
	for _, statement := range fs.Body.Statements {
		if !foundReturnStatement {
			_, foundReturnStatement = statement.(*ast.ReturnStatement)
		}

		c.compileStatement(statement)
	}

	if !foundReturnStatement {
		if _type == VOID {
			c.emit(RETURN)
		} else {
			err := helper.MakeError(fs.Token, fmt.Sprintf("expected return statement in function %q", fs.Var.Name))
			c.addError(err)
		}
	}

	c.emitLabel(end)
}

func (c *Compiler) compileIfStatement(is *ast.IfStatement) {
	elifOrElse := c.getUniqueLabel()
	end := c.getUniqueLabel()

	_type, register := c.compileExpression(is.Condition)

	if _type != builtIn(Bool) {
		err := helper.MakeError(is.Token, fmt.Sprintf("expected boolean condition in if statement, got %q", _type.String()))
		c.addError(err)
	}

	c.emit(COMPARE_WITH_VALUE, register, BOOL_TRUE)
	c.emit(JUMP_IF_NOT_EQUAL, elifOrElse)

	c.enterScope()

	for _, statement := range is.Consequence.Statements {
		c.compileStatement(statement)
	}

	c.emit(JUMP, end)
	c.leaveScope()

	c.emitLabel(elifOrElse)

	if is.Elifs != nil {
		for _, elif := range is.Elifs {
			c.compileElifStatement(elif, end)
		}
	}

	if is.Alternative != nil {
		for _, statement := range is.Alternative.Statements {
			c.compileStatement(statement)
		}
	}

	c.emitLabel(end)
}

func (c *Compiler) compileBreakStatement(bs *ast.BreakStatement) {
	_, begin, ok := c.scope.GetLoopEndAndBegin()
	if !ok {
		err := helper.MakeError(bs.Token, "unexpected Break statement")
		c.addError(err)
	}

	c.emit(JUMP, begin)
}

func (c *Compiler) compileContinueStatement(bs *ast.ContinueStatement) {
	end, _, ok := c.scope.GetLoopEndAndBegin()
	if !ok {
		err := helper.MakeError(bs.Token, "unexpected Continue statement")
		c.addError(err)
	}

	c.emit(JUMP, end)
}

func (c *Compiler) compileElifStatement(es *ast.ElifStatement, end string) {
	nextElif := c.getUniqueLabel()
	_type, register := c.compileExpression(es.Condition)

	if _type != builtIn(Bool) {
		err := helper.MakeError(es.Token, fmt.Sprintf("expected boolean condition in elif statement, got %q", _type.String()))
		c.addError(err)
	}

	c.emit(COMPARE_WITH_VALUE, register, BOOL_TRUE)
	c.emit(JUMP_IF_NOT_EQUAL, nextElif)

	c.enterScope()
	defer c.leaveScope()

	for _, statement := range es.Consequence.Statements {
		c.compileStatement(statement)
	}
	c.emit(JUMP, end)
	c.emitLabel(nextElif)
}

func (c *Compiler) compileExpression(statement ast.Expression) (Type, register) {
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
		return c.compileCallExpression(exp)
	case *ast.ScopeExpression:
		scope, ok := c.findScope(exp, c.scope)
		if !ok {
			err := helper.MakeError(exp.Token, fmt.Sprintf("undeclared scope/alias %q", exp.Scope))
			c.addError(err)
		} else {
			return c.compileIdentifierFromScope(exp.Value, scope)
		}
	default:
		log.Fatalf("type of expression is not handled. got=%T", exp)
	}
	return VOID, ""
}

func (c *Compiler) compileIntegralLiteral(expression *ast.IntegralLiteral) (Type, register) {
	c.emit(LOAD_TO_REG_FROM_VAL, AX, expression.Value)
	return builtIn(Int), AX
}

func (c *Compiler) compileBooleanLiteral(expression *ast.BooleanLiteral) (Type, register) {
	value := BOOL_FALSE
	if expression.Value {
		value = BOOL_TRUE
	}

	c.emit(LOAD_TO_REG_FROM_VAL, AX, value)
	return builtIn(Bool), AX
}

func (c *Compiler) compilePrefixExpression(expression *ast.PrefixExpression) (Type, register) {

	_type, register := c.compileExpression(expression.Right)

	switch expression.Operator {
	case tokens.NOT:
		if _type != builtIn(Bool) {
			err := helper.MakeError(expression.Token, fmt.Sprintf("expected boolean expression. got=%q", _type.String()))
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
	case tokens.NEGATION:
		if _type != builtIn(Int) {
			err := helper.MakeError(expression.Token, fmt.Sprintf("expected integer expression. got=%q", _type.String()))
			c.addError(err)
		}
		log.Fatalf("WIP")
	default:
		log.Fatalf("type of prefix is not handled. got=%q", expression.Operator)
	}

	return builtIn(Bool), register
}

func (c *Compiler) compileInfixExpression(expression *ast.InfixExpression) (Type, register) {

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

	emitComparison := func(jump command) (Type, register) {
		if leftType != builtIn(Int) || rightType != builtIn(Int) {
			err := helper.MakeError(expression.Token, fmt.Sprintf("expected int expression(s). got left=%q and right=%q",
				leftType.String(), rightType.String()))
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
		return builtIn(Bool), AX
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
		if leftType != builtIn(Bool) || rightType != builtIn(Bool) {
			err := helper.MakeError(expression.Token, fmt.Sprintf("expected bool expression(s). got left=%q and right=%q",
				leftType.String(), rightType.String()))
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
		return builtIn(Bool), AX
	case tokens.OR:
		if leftType != builtIn(Bool) || rightType != builtIn(Bool) {
			err := helper.MakeError(expression.Token, fmt.Sprintf("expected bool expression(s). got left=%q and right=%q",
				leftType.String(), rightType.String()))
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
		return builtIn(Bool), AX
	default:
		log.Fatalf("type of infix expression is not handled. got=%q", expression.Operator)
		return VOID, ""
	}
}

func (c *Compiler) compileIdentifier(expression *ast.Identifier) (Type, register) {
	return c.compileIdentifierFromScope(expression, c.scope)
}

func (c *Compiler) compileIdentifierFromScope(expression *ast.Identifier, scope *scope) (Type, register) {
	// TODO check if scope is nil
	if variable, ok := scope.GetVariable(expression.Value); ok {
		c.emit(LOAD_TO_REG_FROM_MEM, AX, variable.Addr)
		return variable.Type, AX
	}

	err := helper.MakeError(expression.Token, fmt.Sprintf("undeclared identifier. got=%q", expression))
	c.addError(err)
	return VOID, ""
}

func (c *Compiler) compileCallExpression(expression *ast.CallExpression) (Type, register) {
	var function name
	var scope *scope

	switch exp := expression.Function.(type) {
	case *ast.ScopeExpression:
		s, ok := c.findScope(exp, c.scope)
		if !ok {
			err := helper.MakeError(exp.Token, fmt.Sprintf("undeclared scope %q", exp.Value.Value))
			c.addError(err)
			return VOID, ""
		}
		function = exp.Value.Value
		scope = s
	case *ast.Identifier:
		function = exp.Value
		scope = c.scope
	default:
		log.Fatalf("type of call expression is not handled. got=%q", expression.Function)
		return VOID, ""
	}

	fun, ok := scope.GetFunction(function)
	if !ok {
		err := helper.MakeError(expression.Token, fmt.Sprintf("undeclared function %q", function))
		c.addError(err)
		return VOID, ""
	}

	if len(fun.Arguments) != len(expression.Arguments) {
		err := helper.MakeError(expression.Token, fmt.Sprintf("unexpected number of arguments expected=%d, got=%d", len(fun.Arguments), len(expression.Arguments)))
		c.addError(err)
		return VOID, ""
	}

	for i := range len(fun.Arguments) {
		arg := fun.Arguments[i]
		passedArg := expression.Arguments[i]
		t, register := c.compileExpression(passedArg)

		if t != arg.Type {
			err := helper.MakeError(expression.Token, fmt.Sprintf("unexpected type of an argument expected %q, got %q", t.String(), arg.Type.String()))
			c.addError(err)
		}

		c.emit(LOAD_TO_MEM_FROM_REG, arg.Addr, register)
	}

	c.emit(CALL, fun.Label)

	return fun.Type, RETURN_REGISTER
}

func (c *Compiler) findScope(expression *ast.ScopeExpression, scope *scope) (*scope, bool) {
	switch exp := expression.Scope.(type) {
	case *ast.ScopeExpression:
		s, ok := c.findScope(exp, scope)
		if !ok {
			err := helper.MakeError(exp.Token, fmt.Sprintf("undeclared scope/alias %q", exp.Value.Value))
			c.addError(err)
		}
		return s.GetScope(exp.Value.Value)
	case *ast.Identifier:
		return scope.GetScope(exp.Value)
	default:
		return scope, false
	}
}

func (c *Compiler) findType(expression *ast.Variable) (Type, bool) {
	switch exp := expression.Type.(type) {
	case *ast.ScopeExpression:
		s, ok := c.findScope(exp, c.scope)
		if !ok {
			err := helper.MakeError(exp.Token, fmt.Sprintf("undeclared scope/alias %q", exp.Value.Value))
			c.addError(err)
			return VOID, false
		}

		return Type{Scope: s, Name: exp.Value.Value}, true
	case *ast.Identifier:
		if slices.Contains(BUILTIN_TYPES, exp.Value) {
			return Type{Scope: nil, Name: exp.Value}, true
		}
		s, ok := c.scope.GetScope(helper.FirstToLowerCase(exp.Value))
		if !ok {
			err := helper.MakeError(exp.Token, fmt.Sprintf("undeclared type %q", exp.Value))
			c.addError(err)
			return VOID, false
		}

		return Type{Scope: s.GetParent(), Name: exp.Value}, true
	default:
		return VOID, false
	}
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
