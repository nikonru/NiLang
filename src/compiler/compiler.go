package compiler

import (
	"NiLang/src/ast"
	"NiLang/src/helper"
	"NiLang/src/lexer"
	"NiLang/src/parser"
	"NiLang/src/tokens"
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type address = int

type Compiler struct {
	output      bytes.Buffer
	memoryIndex address
	variables   []map[name]variable

	labelIndex uint64
}

func New() *Compiler {
	global := make(map[string]variable)
	variables := []map[string]variable{global}
	return &Compiler{memoryIndex: -1, variables: variables, labelIndex: 0}
}

func (c *Compiler) Compile(input []byte) ([]byte, error) {

	lexer := lexer.New(input)
	parser := parser.New(&lexer)
	program := parser.Parse()
	if program == nil {
		log.Fatalf("parser.Parse() has returned nil")
	}

	_errors := parser.Errors()
	if len(_errors) != 0 {
		fmt.Printf("parser had %d error(s)\n", len(_errors))
		fmt.Print("parser error(s):\n")
		for _, err := range _errors {
			helper.PrintError(err, input)
		}

		return c.output.Bytes(), errors.New("parsing errors")
	}

	fmt.Println("PROGRAM TREE")
	for _, statement := range program.Statements {
		c.compileStatement(statement)
		fmt.Println(statement.String())
	}
	fmt.Println("END")
	fmt.Println(c.variables)
	return c.output.Bytes(), nil
}

func (c *Compiler) emit(op command, arg1 interface{}, arg2 interface{}) {
	//TODO make variadic

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
	if arg1 != nil {
		c.output.WriteString(" ")
		write(arg1, 1)

		if arg2 != nil {
			c.output.WriteString(" ")
			write(arg2, 2)
		}
	}

	c.output.WriteString("\n")
}

func (c *Compiler) emitLabel(label string) {
	c.emit(label+":", nil, nil)
}

func (c *Compiler) compileStatement(statement ast.Statement) {
	switch stm := statement.(type) {
	case *ast.DeclarationStatement:
		c.compileDeclarationStatement(stm)
	case *ast.ExpressionStatement:
		c.compileExpression(stm.Expression)
	default:
		log.Fatalf("type of statement is not handled. got=%T", statement)
	}
}

func (c *Compiler) compileDeclarationStatement(ds *ast.DeclarationStatement) {
	_type, register := c.compileExpression(ds.Value)

	if _type != ds.Name.Type.Value {
		log.Fatalf("declared variable and expression have different types. variable=%q, expression=%q", ds.Name.Type.Value, _type)
	}

	addr := c.getMemoryIndex()
	c.emit(LOAD_MEM, addr, register)

	if ok := c.addNewVariable(ds.Name.Value, addr, _type); !ok {
		log.Fatalf("redeclaration of variable %q", ds.Name.Value)
	}
}

func (c *Compiler) compileExpression(statement ast.Expression) (name, register) {
	switch exp := statement.(type) {
	case *ast.IntegralLiteral:
		return c.compileIntegralLiteral(exp)
	case *ast.BooleanLiteral:
		return c.compileBooleanLiteral(exp)
	case *ast.PrefixExpression:
		return c.compilePrefixExpression(exp)
	default:
		log.Fatalf("type of expression is not handled. got=%T", exp)
		return "", ""
	}
}

func (c *Compiler) compileIntegralLiteral(expression *ast.IntegralLiteral) (name, register) {
	c.emit(LOAD_VAL, AX, expression.Value)
	return Int, AX
}

func (c *Compiler) compileBooleanLiteral(expression *ast.BooleanLiteral) (name, register) {
	value := BOOL_FALSE
	if expression.Value {
		value = BOOL_TRUE
	}

	c.emit(LOAD_VAL, AX, value)
	return Bool, AX
}

func (c *Compiler) compilePrefixExpression(expression *ast.PrefixExpression) (name, register) {

	name, register := c.compileExpression(expression.Right)

	switch tokens.LookUpIdent(expression.Operator) {
	case tokens.NOT:
		if name != Bool {
			log.Fatalf("expected boolean expression. got=%q", name)
		}

		end := c.getUniqueLabel()
		True := c.getUniqueLabel()
		False := c.getUniqueLabel()

		c.emit(COMPARE_WITH_VALUE, register, BOOL_TRUE)
		c.emit(JUMP_IF_EQUAL, register, False)

		c.emitLabel(True)
		c.emit(LOAD_VAL, register, BOOL_TRUE)
		c.emit(JUMP, end, nil)

		c.emitLabel(False)
		c.emit(LOAD_VAL, register, BOOL_FALSE)

		c.emitLabel(end)

	default:
		log.Fatalf("type of prefix is not handled. got=%q", expression.Operator)
	}

	return Bool, register
}

func (c *Compiler) getMemoryIndex() address {
	c.memoryIndex++
	return c.memoryIndex
}

func (c *Compiler) getUniqueLabel() string {
	// TODO: maximize number of possible labels
	c.labelIndex++
	return "label" + strconv.FormatUint(c.labelIndex, 10)
}

func (c *Compiler) addNewVariable(name string, addr address, t name) bool {
	scope := c.variables[len(c.variables)-1]

	if _, ok := scope[name]; ok {
		return false
	}

	scope[name] = variable{Addr: addr, Type: t}
	return true
}
