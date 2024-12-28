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

	exp, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		test.Fatalf("expression is not *ast.Identifier, got=%T", statement.Expression)
	}
	if exp.Value != "foobar" {
		test.Errorf("exp.Value is not %s, got=%s", "foobar", exp.Value)
	}
	if exp.TokenLiteral() != "foobar" {
		test.Errorf("exp.TokenLiteral() is not %s, got=%s", "foobar", exp.TokenLiteral())
	}
}

func TestIntegralLiteralExpression(test *testing.T) {
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

	exp, ok := statement.Expression.(*ast.IntegralLiteral)
	if !ok {
		test.Fatalf("expression is not *ast.IntegralLiteral, got=%T", statement.Expression)
	}
	if exp.Value != 5 {
		test.Errorf("exp.Value is not %d, got=%d", 5, exp.Value)
	}
	if exp.TokenLiteral() != "5" {
		test.Errorf("exp.TokenLiteral() is not %s, got=%s", "5", exp.TokenLiteral())
	}
}

func TestTrueBooleanLiteralExpression(test *testing.T) {
	input := []byte(`True`)

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

	exp, ok := statement.Expression.(*ast.BooleanLiteral)
	if !ok {
		test.Fatalf("expression is not *ast.BooleanLiteral, got=%T", statement.Expression)
	}
	if exp.Value != true {
		test.Errorf("exp.Value is not %v, got=%v", true, exp.Value)
	}
	if exp.TokenLiteral() != "True" {
		test.Errorf("exp.TokenLiteral() is not %s, got=%s", "True", exp.TokenLiteral())
	}
}

func TestFalseBooleanLiteralExpression(test *testing.T) {
	input := []byte(`False`)

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

	exp, ok := statement.Expression.(*ast.BooleanLiteral)
	if !ok {
		test.Fatalf("expression is not *ast.BooleanLiteral, got=%T", statement.Expression)
	}
	if exp.Value != false {
		test.Errorf("exp.Value is not %v, got=%v", false, exp.Value)
	}
	if exp.TokenLiteral() != "False" {
		test.Errorf("exp.TokenLiteral() is not %s, got=%s", "False", exp.TokenLiteral())
	}
}

func TestParsingPrefix(test *testing.T) {
	prefixTests := []struct {
		Input    []byte
		Operator string
		Value    bool
	}{
		{[]byte(`Not True`), "Not", true},
		{[]byte(`Not False`), "Not", false},
	}

	for _, testCase := range prefixTests {
		lexer := lexer.New(testCase.Input)
		parser := parser.New(&lexer)
		program := parser.Parse()
		checkParseErrors(test, parser, testCase.Input)

		if len(program.Statements) != 1 {
			test.Fatalf("program has not enough statements, got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			test.Fatalf("expression is not *ast.PrefixExpression, got=%T", statement.Expression)
		}
		if exp.Operator != testCase.Operator {
			test.Errorf("exp.operator is not %v, got=%v", testCase.Operator, exp.Operator)
		}
		if !testBooleanLiteral(test, exp.Right, testCase.Value) {
			return
		}
	}
}

func testBooleanLiteral(test *testing.T, expression ast.Expression, value bool) bool {
	boolean, ok := expression.(*ast.BooleanLiteral)
	if !ok {
		test.Errorf("expression is not *ast.BooleanLiteral, got=%T", expression)
		return false
	}

	if boolean.Value != value {
		test.Errorf("expression is not %v, got=%v", value, expression)
		return false
	}

	if boolean.TokenLiteral() != "True" && boolean.TokenLiteral() != "False" {
		test.Errorf("boolean.TokenLiteral()  is not True or False, got=%s", boolean.TokenLiteral())
		return false
	}

	return true
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
