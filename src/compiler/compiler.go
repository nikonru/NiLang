package compiler

import (
	"NiLang/src/helper"
	"NiLang/src/lexer"
	"NiLang/src/parser"
	"NiLang/src/ast"
	"bytes"
	"errors"
	"strconv"
	"fmt"
	"log"
)

type Compiler struct {
	output bytes.Buffer
}

func New() *Compiler {
	return &Compiler{}
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
		fmt.Println(statement.String())
	}
	fmt.Println("END")

	c.emit(LOAD, AX, BX)
	c.emit(LOAD, 2, BX)
	c.emitLabel("lab")

	return c.output.Bytes(), nil
}

func (c *Compiler) emit(op command, arg1 interface{}, arg2 interface{}) {

	write := func(arg interface{}, id int){
		switch v := arg.(type) {
		case int:
			c.output.WriteString(strconv.Itoa(v))
		case string:
			c.output.WriteString(v)
		case bool:
			if v{
				c.output.WriteString("1")
			} else{
				c.output.WriteString("0")
			}
		default:
			log.Fatalf("type of arg%d not handled. got=%T", arg, id)
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

func (c *Compiler) emitLabel(label string){
	c.emit(label +":", nil, nil)
}

func (c *Compiler) compileDeclarationStatement(ds *ast.DeclarationStatement){
	return
}

