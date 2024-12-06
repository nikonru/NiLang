package compiler

import (
	"NiLang/src/lexer"
	"NiLang/src/tokens"
	"fmt"
)

type Compiler struct {
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(input []byte) ([]byte, error) {

	Lexer := lexer.New(input)

	for token := Lexer.NextToken(); token.Type != tokens.EOF; token = Lexer.NextToken() {
		fmt.Println(token)
	}

	return input, nil
}
