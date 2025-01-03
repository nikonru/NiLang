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

func TestUsingStatement(test *testing.T) {
	input := []byte(`Using bot
Using world`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)

	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	length := 2
	if len(program.Statements) != length {
		test.Fatalf("program.Statements doesn't contain %d statements: got=%v", length, len(program.Statements))
	}

	tests := []struct {
		expectedTypeLiteral string
		expectedIdentifier  string
	}{
		{"Using", "bot"},
		{"Using", "world"},
	}

	for i, t := range tests {
		statement := program.Statements[i]
		if !testUsingStatement(test, statement, t.expectedTypeLiteral, t.expectedIdentifier) {
			return
		}
	}
}

func TestIdentifierExpression(test *testing.T) {
	input := []byte(`foobar`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)
	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
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
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
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
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
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
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
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

func TestParsingPrefixExpression(test *testing.T) {
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
		if program == nil {
			test.Fatalf("parser.Parse() has returned nil")
		}
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
			test.Errorf("exp.Operator is not %v, got=%v", testCase.Operator, exp.Operator)
		}
		if !testBooleanLiteral(test, exp.Right, testCase.Value) {
			return
		}
	}
}

func TestParsingInfixExpression(test *testing.T) {
	infixTests := []struct {
		Input      []byte
		LeftValue  int64
		Operator   string
		RightValue int64
	}{
		{[]byte(`5 > 6`), 5, ">", 6},
		{[]byte(`5 >= 6`), 5, ">=", 6},
		{[]byte(`6 < 7`), 6, "<", 7},
		{[]byte(`6 <= 7`), 6, "<=", 7},
		{[]byte(`6 == 6`), 6, "==", 6},
		{[]byte(`1 != 20`), 1, "!=", 20},
	}

	for _, testCase := range infixTests {
		lexer := lexer.New(testCase.Input)
		parser := parser.New(&lexer)
		program := parser.Parse()
		if program == nil {
			test.Fatalf("parser.Parse() has returned nil")
		}
		checkParseErrors(test, parser, testCase.Input)

		if len(program.Statements) != 1 {
			test.Fatalf("program has not enough statements, got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			test.Fatalf("expression is not *ast.InfixExpression, got=%T", statement.Expression)
		}

		if !testInfixExpression(test, exp, testCase.LeftValue, testCase.Operator, testCase.RightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(test *testing.T) {
	tests := []struct {
		Input    []byte
		Expected string
	}{
		{[]byte(`5 < 6 == Not True`), "((5 < 6) == (NotTrue))"},
		{[]byte(`5 >= 6 <= 10`), "((5 >= 6) <= 10)"},
		{[]byte(`Not  True   ==  False  `), "((NotTrue) == False)"},
	}

	for _, testCase := range tests {
		lexer := lexer.New(testCase.Input)
		parser := parser.New(&lexer)
		program := parser.Parse()
		if program == nil {
			test.Fatalf("parser.Parse() has returned nil")
		}
		checkParseErrors(test, parser, testCase.Input)

		actual := program.String()
		if actual != testCase.Expected {
			test.Errorf("expected=%q, got=%q", testCase.Expected, actual)
		}
	}
}

func TestIfExpression(test *testing.T) {
	input := []byte(`
If x > 5:
    Foo
Elif z > 3:
    Bar
Elif z < x:
    Car
Else:
    Goo
    `)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)
	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	if len(program.Statements) != 1 {
		test.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		test.Fatalf("statement.Expression is not ast.IfExpression, got=%T", statement.Expression)
	}

	condition, ok := exp.Condition.(*ast.InfixExpression)
	if !ok {
		test.Fatalf("exp.Condition is not InfixExpression, got=%T", exp.Condition)
	}

	if !testInfixExpression(test, condition, "x", ">", 5) {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		test.Fatalf("consequence is not 1 statements, got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(test, consequence.Expression, "Foo") {
		return
	}

	if len(exp.Elifs) != 2 {
		test.Fatalf("elifs is not 2 statements, got=%d", len(exp.Elifs))
	}

	elifCondition1, ok := exp.Elifs[0].Condition.(*ast.InfixExpression)
	if !ok {
		test.Fatalf("exp.Elifs[0].Condition is not ast.InfixExpression, got=%T", exp.Elifs[0].Condition)
	}

	if !testInfixExpression(test, elifCondition1, "z", ">", 3) {
		return
	}

	if len(exp.Elifs[0].Consequence.Statements) != 1 {
		test.Fatalf("exp.Elifs[0].Consequence.Statements is not 1 statements, got=%d", len(exp.Elifs[0].Consequence.Statements))
	}

	consequence1, ok := exp.Elifs[0].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("exp.Elifs[0].Consequence.Statements[0] is not ast.ExpressionStatement, got=%T", exp.Elifs[0].Consequence.Statements[0])
	}

	if !testIdentifier(test, consequence1.Expression, "Bar") {
		return
	}

	elifCondition2, ok := exp.Elifs[1].Condition.(*ast.InfixExpression)
	if !ok {
		test.Fatalf("exp.Elifs[1].Condition is not ast.InfixExpression, got=%T", exp.Elifs[1].Condition)
	}

	if !testInfixExpression(test, elifCondition2, "z", "<", "x") {
		return
	}

	if len(exp.Elifs[1].Consequence.Statements) != 1 {
		test.Fatalf("exp.Elifs[1].Consequence.Statements is not 1 statements, got=%d", len(exp.Elifs[1].Consequence.Statements))
	}

	consequence2, ok := exp.Elifs[1].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("exp.Elifs[1].Consequence.Statements[0] is not ast.ExpressionStatement, got=%T", exp.Elifs[1].Consequence.Statements[0])
	}

	if !testIdentifier(test, consequence2.Expression, "Car") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		test.Fatalf("alternative is not 1 statements, got=%d", len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement, got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(test, alternative.Expression, "Goo") {
		return
	}
}

func testInfixExpression(test *testing.T, exp *ast.InfixExpression, left interface{}, operator string, right interface{}) bool {

	if !testLiteralExpression(test, exp.Left, left) {
		return false
	}
	if exp.Operator != operator {
		test.Errorf("exp.Operator is not %v, got=%v", operator, exp.Operator)
		return false
	}
	if !testLiteralExpression(test, exp.Right, right) {
		return false
	}
	return true
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

func testIntegralLiteral(test *testing.T, expression ast.Expression, value int64) bool {
	integral, ok := expression.(*ast.IntegralLiteral)
	if !ok {
		test.Errorf("expression is not *ast.IntegralLiteral, got=%T", expression)
		return false
	}

	if integral.Value != value {
		test.Errorf("expression is not %v, got=%v", value, expression)
		return false
	}

	if integral.TokenLiteral() != fmt.Sprintf("%d", value) {
		test.Errorf("boolean.TokenLiteral() is not %d, got=%s", value, integral.TokenLiteral())
		return false
	}

	return true
}

func testUsingStatement(test *testing.T, statement ast.Statement, literal string, name string) bool {
	if statement.TokenLiteral() != literal {
		test.Errorf("statement.TokenLiteral() is not Bool: got=%v", statement.TokenLiteral())
		return false
	}

	usingStatement, ok := statement.(*ast.UsingStatement)
	if !ok {
		test.Errorf("statement is not *ast.UsingStatement type: got=%v", statement)
		return false
	}

	if usingStatement.Name.Value != name {
		test.Errorf("usingStatement.Name.Value is not *%v type: got=%v", name, statement)
		return false
	}

	if usingStatement.Name.TokenLiteral() != name {
		test.Errorf("usingStatement.Name.TokenLiteral is not *%v type: got=%v", name, statement)
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

	test.Errorf("parser had %d error(s)", len(errors))
	fmt.Print("parser error(s):\n")
	for _, err := range errors {
		helper.PrintError(err, input)
	}
	test.FailNow()
}

func testLiteralExpression(test *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegralLiteral(test, exp, int64(v))
	case int64:
		return testIntegralLiteral(test, exp, v)
	case string:
		return testIdentifier(test, exp, v)
	default:
		test.Errorf("type of exp not handled. got=%T", exp)
		return false
	}
}

func testIdentifier(test *testing.T, expression ast.Expression, value string) bool {
	ident, ok := expression.(*ast.Identifier)
	if !ok {
		test.Errorf("expression is not *ast.Identifier, got=%T", expression)
		return false
	}

	if ident.Value != value {
		test.Errorf("expression is not %v, got=%v", value, expression)
		return false
	}

	if ident.TokenLiteral() != value {
		test.Errorf("ident.TokenLiteral() is not %s, got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}
