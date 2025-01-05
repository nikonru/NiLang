package compiler

import (
	"NiLang/src/helper"
	"NiLang/src/lexer"
	"NiLang/src/parser"
	"fmt"
	"log"
)

type Compiler struct {
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

	errors := parser.Errors()
	if len(errors) != 0 {
		fmt.Printf("parser had %d error(s)\n", len(errors))
		fmt.Print("parser error(s):\n")
		for _, err := range errors {
			helper.PrintError(err, input)
		}
	}

	fmt.Println("PROGRAM TREE")
	for _, statement := range program.Statements {
		fmt.Println(statement.String())
	}
	fmt.Println("END")
	return input, nil
}
