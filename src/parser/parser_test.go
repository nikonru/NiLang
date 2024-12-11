package parser_test

import (
	"NiLang/src/ast"
	"NiLang/src/lexer"
	"NiLang/src/parser"
	"testing"
)

func TestLetStatement(test *testing.T) {
	input := []byte(`
    Bool x = false
    Int number = 1200
    `)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)

	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}

	if len(program.Statements) != 3 {
		test.Fatalf("program.Statements doesn't contain 2 statements: got=%v", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"number"},
	}

	for i, t := range tests {
		statement := program.Statements[i]
		if !testLetStatement(test, statement, t.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(test *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "Bool" {
		test.Errorf("statement.TokenLiteral() is not Bool: got=%v", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		test.Errorf("statement is not *ast.LetStatement type: got=%v", statement)
		return false
	}

	if letStatement.Name.Value != name {
		test.Errorf("letStatement.Name.Value is not *%v type: got=%v", name, statement)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		test.Errorf("letStatement.Name.TokenLiteral is not *%v type: got=%v", name, statement)
		return false
	}

	return true
}
