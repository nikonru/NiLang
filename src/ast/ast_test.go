package ast_test

import (
	"NiLang/src/ast"
	"NiLang/src/tokens"
	"testing"
)

func TestString(test *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.DeclarationStatement{
				Token: tokens.Token{Type: tokens.INT, Literal: "Int", Line: 1, Offset: 0},
				Name:  &ast.Identifier{Token: tokens.Token{Type: tokens.IDENT, Literal: "myVar", Line: 1, Offset: 4}, Value: "myVar"},
				Value: &ast.Identifier{Token: tokens.Token{Type: tokens.IDENT, Literal: "anotherVar", Line: 1, Offset: 9}, Value: "anotherVar"},
			},
		},
	}

	if program.String() != "Int myVar = anotherVar" {
		test.Errorf("program.String() wrong, got=%q", program.String())
	}
}
