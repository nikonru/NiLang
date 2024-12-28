package parser_test

import (
	"NiLang/src/ast"
	"NiLang/src/helper"
	"NiLang/src/lexer"
	"NiLang/src/parser"
	"fmt"
	"testing"
)

func TestDeclarationStatement(test *testing.T) {
	input := []byte(`Bool x = false
Int number = 1400
Dir face = forward`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)

	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	length := 3
	if len(program.Statements) != length {
		test.Fatalf("program.Statements doesn't contain %d statements: got=%v", length, len(program.Statements))
	}

	tests := []struct {
		expectedTypeLiteral string
		expectedIdentifier  string
	}{
		{"Bool", "x"},
		{"Int", "number"},
		{"Dir", "face"},
	}

	for i, t := range tests {
		statement := program.Statements[i]
		if !testDeclarationStatement(test, statement, t.expectedTypeLiteral, t.expectedIdentifier) {
			return
		}
	}
}

func TestIdentifierExpression(test *testing.T) {
	input := []byte(`foobar`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)
	program := parser.Parse()
	checkParseErrors(test, parser, input)

	if len(program.Statements) != 1 {
		test.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		test.Fatalf("expression not *ast.Identifier, got=%T", statement.Expression)
	}
	if ident.Value != "foobar" {
		test.Errorf("ident.Value not %s, got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		test.Errorf("ident.TokenLiteral() not %s, got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegralLiteralExpressio(test *testing.T) {
	input := []byte(`5`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)
	program := parser.Parse()
	checkParseErrors(test, parser, input)

	if len(program.Statements) != 1 {
		test.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.IntegralLiteral)
	if !ok {
		test.Fatalf("expression not *ast.IntegralLiteral, got=%T", statement.Expression)
	}
	if ident.Value != 5 {
		test.Errorf("ident.Value not %d, got=%d", 5, ident.Value)
	}
	if ident.TokenLiteral() != "5" {
		test.Errorf("ident.TokenLiteral() not %s, got=%s", "5", ident.TokenLiteral())
	}
}

func testDeclarationStatement(test *testing.T, statement ast.Statement, literal string, name string) bool {
	if statement.TokenLiteral() != literal {
		test.Errorf("statement.TokenLiteral() is not Bool: got=%v", statement.TokenLiteral())
		return false
	}

	declarationStatement, ok := statement.(*ast.DeclarationStatement)
	if !ok {
		test.Errorf("statement is not *ast.DeclarationStatement type: got=%v", statement)
		return false
	}

	if declarationStatement.Name.Value != name {
		test.Errorf("declarationStatement.Name.Value is not *%v type: got=%v", name, statement)
		return false
	}

	if declarationStatement.Name.TokenLiteral() != name {
		test.Errorf("declarationStatement.Name.TokenLiteral is not *%v type: got=%v", name, statement)
		return false
	}

	return true
}

func checkParseErrors(test *testing.T, parser *parser.Parser, input []byte) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}

	test.Errorf("parser had %d errors", len(errors))
	fmt.Print("parser error:\n")
	for _, err := range errors {
		helper.PrintError(err, input)
	}
	test.FailNow()
}
