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
	// TODO check expression
	for i, t := range tests {
		statement := program.Statements[i]
		if !testDeclarationStatement(test, statement, t.expectedTypeLiteral, t.expectedIdentifier) {
			return
		}
	}
}

func TestAssignmentStatement(test *testing.T) {
	input := []byte(`
x = 100
hungry = False`)

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
		Name  string
		Value interface{}
	}{
		{"x", 100},
		{"hungry", false},
	}

	for i, t := range tests {
		statement, ok := program.Statements[i].(*ast.AssignmentStatement)
		if !ok {
			test.Fatalf("statement is not *ast.AssignmentStatement type: got=%v", statement)
		}

		if !testIdentifier(test, statement.Name, t.Name) {
			return
		}

		if !testLiteralExpression(test, statement.Value, t.Value) {
			return
		}
	}
}

func TestAliasStatement(test *testing.T) {
	input := []byte(`
Alias Numbers::Int:
    one = 1
    two = 2
    four = 4`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)

	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	length := 1
	if len(program.Statements) != length {
		test.Fatalf("program.Statements doesn't contain %d statements: got=%v", length, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.AliasStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ReturnStatement, got=%T", program.Statements[0])
	}

	expectedType := "Int"
	if !testTypedIdentifier(test, statement.Name, expectedType, "Numbers") {
		return
	}

	tests := []string{"one", "two", "four"}

	// TODO check expression
	for i, ident := range tests {
		if !testDeclarationStatement(test, statement.Values[i], expectedType, ident) {
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

func TestReturnStatement(test *testing.T) {
	input := []byte(`Return 2`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)

	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	length := 1
	if len(program.Statements) != length {
		test.Fatalf("program.Statements doesn't contain %d statements: got=%v", length, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ReturnStatement, got=%T", program.Statements[0])
	}

	if !testIntegralLiteral(test, statement.Value, 2) {
		return
	}
}

func TestScopeStatement(test *testing.T) {
	input := []byte(`
Scope names:
    Int fish = 1
    Bool flag = False`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)

	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	length := 1
	if len(program.Statements) != length {
		test.Fatalf("program.Statements doesn't contain %d statements: got=%v", length, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ScopeStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ScopeStatement, got=%T", program.Statements[0])
	}
	// TODO check expression

	length = 2
	if len(statement.Body.Statements) != length {
		test.Fatalf("wrong number of statements in body of the scope, expected=%d, got=%T", length, len(statement.Body.Statements))
	}

	if !testDeclarationStatement(test, statement.Body.Statements[0], "Int", "fish") {
		return
	}

	if !testDeclarationStatement(test, statement.Body.Statements[1], "Bool", "flag") {
		return
	}
}

func TestWhileStatement(test *testing.T) {
	input := []byte(`
While 1 < 2:
    Foo
`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)

	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	length := 1
	if len(program.Statements) != length {
		test.Fatalf("program.Statements doesn't contain %d statements: got=%v", length, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.WhileStatement, got=%T", program.Statements[0])
	}

	condition, ok := statement.Condition.(*ast.InfixExpression)
	if !ok {
		test.Fatalf("statement.Condition is not ast.WhileStatement, got=%T", statement.Condition)
	}

	if !testInfixExpression(test, condition, 1, "<", 2) {
		return
	}

	length = 1
	if len(statement.Body.Statements) != length {
		test.Fatalf("wrong number of statements in body of the scope, expected=%d, got=%T", length, len(statement.Body.Statements))
	}

	exp, ok := statement.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("statement.Body.Statements[0] is not ast.ScopeStatement, got=%T", statement.Body.Statements[0])
	}

	if !testEmptyCallExpression(test, exp.Expression, "Foo") {
		return
	}
}

func TestFunctionStatement(test *testing.T) {
	input := []byte(`
Fun F::Bool$max Int, default Bool:
    Return 5 > max

Fun I::Int$v Int:
    Return v

Fun G::Int:
    Foo
    Return 5

Fun H:
    Foo

Fun Z$v Int:
    Foo
`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)

	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	length := 5
	if len(program.Statements) != length {
		test.Fatalf("program.Statements doesn't contain %d statements: got=%v", length, len(program.Statements))
	}

	//-0-

	statement, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.FunnctionStatement, got=%T", program.Statements[0])
	}

	if statement.TokenLiteral() != "Fun" {
		test.Fatalf("statement.TokenLiteral() is not `Fun`, got=%q", statement.TokenLiteral())
	}

	if !testTypedIdentifier(test, statement.Name, "Bool", "F") {
		return
	}

	length = 2
	if len(statement.Parameters) != length {
		test.Fatalf("wrong number of parameters in F, expected=%d, got=%d", length, len(statement.Parameters))
	}

	if !testTypedIdentifier(test, statement.Parameters[0], "Int", "max") {
		return
	}

	if !testTypedIdentifier(test, statement.Parameters[1], "Bool", "default") {
		return
	}

	length = 1
	if len(statement.Body.Statements) != length {
		test.Fatalf("wrong number of statements in F, expected=%d, got=%d", length, len(statement.Body.Statements))
	}

	ret, ok := statement.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		test.Fatalf("statement.Body.Statements[0] is not ast.ReturnStatement, got=%T", statement.Body.Statements[0])
	}

	if ret.TokenLiteral() != "Return" {
		test.Fatalf("ret.TokenLiteral() is not Return, got=%q", ret.TokenLiteral())
	}

	exp, ok := ret.Value.(*ast.InfixExpression)
	if !ok {
		test.Fatalf("ret.Value is not ast.InfixExpression, got=%T", ret.Value)
	}

	if !testInfixExpression(test, exp, 5, ">", "max") {
		return
	}

	//-1-

	statement, ok = program.Statements[1].(*ast.FunctionStatement)
	if !ok {
		test.Fatalf("program.Statements[1] is not ast.FunnctionStatement, got=%T", program.Statements[1])
	}

	if statement.TokenLiteral() != "Fun" {
		test.Fatalf("statement.TokenLiteral() is not `Fun`, got=%q", statement.TokenLiteral())
	}

	if !testTypedIdentifier(test, statement.Name, "Int", "I") {
		return
	}

	length = 1
	if len(statement.Parameters) != length {
		test.Fatalf("wrong number of parameters in I, expected=%d, got=%d", length, len(statement.Parameters))
	}

	if !testTypedIdentifier(test, statement.Parameters[0], "Int", "v") {
		return
	}

	length = 1
	if len(statement.Body.Statements) != length {
		test.Fatalf("wrong number of statements in I, expected=%d, got=%d", length, len(statement.Body.Statements))
	}

	ret, ok = statement.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		test.Fatalf("statement.Body.Statements[0] is not ast.ReturnStatement, got=%T", statement.Body.Statements[0])
	}

	if ret.TokenLiteral() != "Return" {
		test.Fatalf("ret.TokenLiteral() is not Return, got=%q", ret.TokenLiteral())
	}

	ident, ok := ret.Value.(*ast.Identifier)
	if !ok {
		test.Fatalf("ret.Value is not ast.InfixExpression, got=%T", ret.Value)
	}

	if !testIdentifier(test, ident, "v") {
		return
	}

	//-2-

	statement, ok = program.Statements[2].(*ast.FunctionStatement)
	if !ok {
		test.Fatalf("program.Statements[2] is not ast.FunnctionStatement, got=%T", program.Statements[2])
	}

	if statement.TokenLiteral() != "Fun" {
		test.Fatalf("statement.TokenLiteral() is not `Fun`, got=%q", statement.TokenLiteral())
	}

	if !testTypedIdentifier(test, statement.Name, "Int", "G") {
		return
	}

	length = 0
	if len(statement.Parameters) != length {
		test.Fatalf("wrong number of parameters in G, expected=%d, got=%d", length, len(statement.Parameters))
	}

	length = 2
	if len(statement.Body.Statements) != length {
		test.Fatalf("wrong number of statements in G, expected=%d, got=%d", length, len(statement.Body.Statements))
	}

	expStatement, ok := statement.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("statement.Body.Statements[0] is not ast.ExpressionStatement, got=%T", statement.Body.Statements[0])
	}

	if !testEmptyCallExpression(test, expStatement.Expression, "Foo") {
		return
	}

	ret, ok = statement.Body.Statements[1].(*ast.ReturnStatement)
	if !ok {
		test.Fatalf("statement.Body.Statements[1] is not ast.ReturnStatement, got=%T", statement.Body.Statements[0])
	}

	if ret.TokenLiteral() != "Return" {
		test.Fatalf("ret.TokenLiteral() is not Return, got=%q", ret.TokenLiteral())
	}

	integer, ok := ret.Value.(*ast.IntegralLiteral)
	if !ok {
		test.Fatalf("ret.Value is not ast.InfixExpression, got=%T", ret.Value)
	}

	if !testIntegralLiteral(test, integer, 5) {
		return
	}

	//-3-

	statement, ok = program.Statements[3].(*ast.FunctionStatement)
	if !ok {
		test.Fatalf("program.Statements[3] is not ast.FunnctionStatement, got=%T", program.Statements[3])
	}

	if statement.TokenLiteral() != "Fun" {
		test.Fatalf("statement.TokenLiteral() is not `Fun`, got=%q", statement.TokenLiteral())
	}

	if statement.Name.Type != nil {
		test.Fatalf("statement.Name.Type is not `nil`, got=%v", statement.Name.Type)
	}

	if statement.Name.TokenLiteral() != "H" {
		test.Fatalf("statement.Name.TokenLiteral() is not '%s', got=%s", "H", statement.Name.TokenLiteral())
	}

	if statement.Name.Value != "H" {
		test.Fatalf("statement.Name.Value is not '%s', got=%s", "H", statement.Name.Value)
	}

	length = 0
	if len(statement.Parameters) != length {
		test.Fatalf("wrong number of parameters in H, expected=%d, got=%d", length, len(statement.Parameters))
	}

	length = 1
	if len(statement.Body.Statements) != length {
		test.Fatalf("wrong number of statements in H, expected=%d, got=%d", length, len(statement.Body.Statements))
	}

	expStatement, ok = statement.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("statement.Body.Statements[0] is not ast.ExpressionStatement, got=%T", statement.Body.Statements[0])
	}

	if !testEmptyCallExpression(test, expStatement.Expression, "Foo") {
		return
	}

	//-4-

	statement, ok = program.Statements[4].(*ast.FunctionStatement)
	if !ok {
		test.Fatalf("program.Statements[4] is not ast.FunnctionStatement, got=%T", program.Statements[4])
	}

	if statement.TokenLiteral() != "Fun" {
		test.Fatalf("statement.TokenLiteral() is not `Fun`, got=%q", statement.TokenLiteral())
	}

	if statement.Name.Type != nil {
		test.Fatalf("statement.Name.Type is not `nil`, got=%v", statement.Name.Type)
	}

	if statement.Name.TokenLiteral() != "Z" {
		test.Fatalf("statement.Name.TokenLiteral() is not '%s', got=%s", "Z", statement.Name.TokenLiteral())
	}

	if statement.Name.Value != "Z" {
		test.Fatalf("statement.Name.Value is not '%s', got=%s", "Z", statement.Name.Value)
	}

	length = 1
	if len(statement.Parameters) != length {
		test.Fatalf("wrong number of parameters in H, expected=%d, got=%d", length, len(statement.Parameters))
	}

	length = 1
	if len(statement.Body.Statements) != length {
		test.Fatalf("wrong number of statements in H, expected=%d, got=%d", length, len(statement.Body.Statements))
	}

	if !testTypedIdentifier(test, statement.Parameters[0], "Int", "v") {
		return
	}

	expStatement, ok = statement.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("statement.Body.Statements[0] is not ast.ExpressionStatement, got=%T", statement.Body.Statements[0])
	}

	if !testEmptyCallExpression(test, expStatement.Expression, "Foo") {
		return
	}
}

func TestIdentifierExpression(test *testing.T) {
	input := []byte(`
foobar
goo`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)
	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	length := 2
	if len(program.Statements) != 2 {
		test.Fatalf("program has not enough statements, expectdd=%d, got=%d", length, len(program.Statements))
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

	statement, ok = program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[1] is not ast.ExpressionStatement, got=%T", program.Statements[1])
	}

	exp, ok = statement.Expression.(*ast.Identifier)
	if !ok {
		test.Fatalf("expression is not *ast.Identifier, got=%T", statement.Expression)
	}
	if exp.Value != "goo" {
		test.Errorf("exp.Value is not %s, got=%s", "goo", exp.Value)
	}
	if exp.TokenLiteral() != "goo" {
		test.Errorf("exp.TokenLiteral() is not %s, got=%s", "goo", exp.TokenLiteral())
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

func TestIfStatement(test *testing.T) {
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

	statement, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.IfExpression, got=%T", program.Statements[0])
	}

	condition, ok := statement.Condition.(*ast.InfixExpression)
	if !ok {
		test.Fatalf("exp.Condition is not InfixExpression, got=%T", statement.Condition)
	}

	if !testInfixExpression(test, condition, "x", ">", 5) {
		return
	}

	if len(statement.Consequence.Statements) != 1 {
		test.Fatalf("consequence is not 1 statements, got=%d", len(statement.Consequence.Statements))
	}

	consequence, ok := statement.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement, got=%T", statement.Consequence.Statements[0])
	}

	if !testEmptyCallExpression(test, consequence.Expression, "Foo") {
		return
	}

	if len(statement.Elifs) != 2 {
		test.Fatalf("elifs is not 2 statements, got=%d", len(statement.Elifs))
	}

	elifCondition1, ok := statement.Elifs[0].Condition.(*ast.InfixExpression)
	if !ok {
		test.Fatalf("exp.Elifs[0].Condition is not ast.InfixExpression, got=%T", statement.Elifs[0].Condition)
	}

	if !testInfixExpression(test, elifCondition1, "z", ">", 3) {
		return
	}

	if len(statement.Elifs[0].Consequence.Statements) != 1 {
		test.Fatalf("exp.Elifs[0].Consequence.Statements is not 1 statements, got=%d", len(statement.Elifs[0].Consequence.Statements))
	}

	consequence1, ok := statement.Elifs[0].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("exp.Elifs[0].Consequence.Statements[0] is not ast.ExpressionStatement, got=%T", statement.Elifs[0].Consequence.Statements[0])
	}

	if !testEmptyCallExpression(test, consequence1.Expression, "Bar") {
		return
	}

	elifCondition2, ok := statement.Elifs[1].Condition.(*ast.InfixExpression)
	if !ok {
		test.Fatalf("exp.Elifs[1].Condition is not ast.InfixExpression, got=%T", statement.Elifs[1].Condition)
	}

	if !testInfixExpression(test, elifCondition2, "z", "<", "x") {
		return
	}

	if len(statement.Elifs[1].Consequence.Statements) != 1 {
		test.Fatalf("exp.Elifs[1].Consequence.Statements is not 1 statements, got=%d", len(statement.Elifs[1].Consequence.Statements))
	}

	consequence2, ok := statement.Elifs[1].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("exp.Elifs[1].Consequence.Statements[0] is not ast.ExpressionStatement, got=%T", statement.Elifs[1].Consequence.Statements[0])
	}

	if !testEmptyCallExpression(test, consequence2.Expression, "Car") {
		return
	}

	if len(statement.Alternative.Statements) != 1 {
		test.Fatalf("alternative is not 1 statements, got=%d", len(statement.Alternative.Statements))
	}

	alternative, ok := statement.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement, got=%T", statement.Alternative.Statements[0])
	}

	if !testEmptyCallExpression(test, alternative.Expression, "Goo") {
		return
	}
}

func TestCallExpression(test *testing.T) {
	input := []byte(`
Get$1 < x,z,  z==2
Consume
Read$2`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)
	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	length := 3
	if len(program.Statements) != length {
		test.Fatalf("program has not enough statements, expected=%d, got=%d", length, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		test.Fatalf("statement.Expression is not ast.CallExpression, got=%T", statement.Expression)
	}

	if !testIdentifier(test, exp.Function, "Get") {
		return
	}

	length = 3
	if len(exp.Arguments) != length {
		test.Fatalf("wrong length of arguments, expected=%d, got=%d", length, len(exp.Arguments))
	}

	arg1, ok := exp.Arguments[0].(*ast.InfixExpression)
	if !ok {
		test.Fatalf("exp.Arguments[0] is not ast.InfixExpression, got=%T", exp.Arguments[0])
	}

	if !testInfixExpression(test, arg1, 1, "<", "x") {
		return
	}

	if !testLiteralExpression(test, exp.Arguments[1], "z") {
		return
	}

	arg3, ok := exp.Arguments[2].(*ast.InfixExpression)
	if !ok {
		test.Fatalf("exp.Arguments[2] is not ast.InfixExpression, got=%T", exp.Arguments[2])
	}

	if !testInfixExpression(test, arg3, "z", "==", 2) {
		return
	}

	//

	statement, ok = program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[1] is not ast.ExpressionStatement, got=%T", program.Statements[1])
	}

	exp, ok = statement.Expression.(*ast.CallExpression)
	if !ok {
		test.Fatalf("statement.Expression is not ast.CallExpression, got=%T", statement.Expression)
	}

	if !testIdentifier(test, exp.Function, "Consume") {
		return
	}

	length = 0
	if len(exp.Arguments) != length {
		test.Fatalf("wrong length of arguments, expected=%d, got=%d", length, len(exp.Arguments))
	}

	//

	statement, ok = program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[2] is not ast.ExpressionStatement, got=%T", program.Statements[2])
	}

	exp, ok = statement.Expression.(*ast.CallExpression)
	if !ok {
		test.Fatalf("statement.Expression is not ast.CallExpression, got=%T", statement.Expression)
	}

	if !testIdentifier(test, exp.Function, "Read") {
		return
	}

	length = 1
	if len(exp.Arguments) != length {
		test.Fatalf("wrong length of arguments, expected=%d, got=%d", length, len(exp.Arguments))
	}

	if !testIntegralLiteral(test, exp.Arguments[0], 2) {
		return
	}
}

func TestScopeExpression(test *testing.T) {
	input := []byte(`bot::max
world::Get
world::compass::south
world::Set$1`)

	lexer := lexer.New(input)
	parser := parser.New(&lexer)
	program := parser.Parse()
	if program == nil {
		test.Fatalf("parser.Parse() has returned nil")
	}
	checkParseErrors(test, parser, input)

	length := 4
	if len(program.Statements) != length {
		test.Fatalf("program has not enough statements, expected=%d, got=%d", length, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := statement.Expression.(*ast.ScopeExpression)
	if !ok {
		test.Fatalf("statement.Expression is not ast.ScopeExpression, got=%T", statement.Expression)
	}

	scope, ok := exp.Scope.(*ast.Identifier)
	if !ok {
		test.Fatalf("exp.Scope is not ast.Identifier, got=%T", exp.Scope)
	}

	if !testIdentifier(test, scope, "bot") {
		return
	}

	ident, ok := exp.Value.(*ast.Identifier)
	if !ok {
		test.Fatalf("exp.Value is not ast.Identifier, got=%T", exp.Value)
	}

	if !testIdentifier(test, ident, "max") {
		return
	}

	//

	statement, ok = program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[1] is not ast.ExpressionStatement, got=%T", program.Statements[1])
	}

	call, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		test.Fatalf("statement.Expression is not ast.CallExpression, got=%T", statement.Expression)
	}

	if len(call.Arguments) != 0 {
		test.Fatalf("wrong number of arguments,expected=%d, got=%T", 0, call.TokenLiteral())
	}

	superScope, ok := call.Function.(*ast.ScopeExpression)
	if !ok {
		test.Fatalf("call.Function is not ast.ScopeExpression, got=%T", exp.Scope)
	}

	if !testIdentifier(test, superScope.Scope, "world") {
		return
	}

	if !testIdentifier(test, superScope.Value, "Get") {
		return
	}

	//

	statement, ok = program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[2] is not ast.ExpressionStatement, got=%T", program.Statements[2])
	}

	exp, ok = statement.Expression.(*ast.ScopeExpression)
	if !ok {
		test.Fatalf("statement.Expression is not ast.ScopeExpression, got=%T", statement.Expression)
	}

	superScope, ok = exp.Scope.(*ast.ScopeExpression)
	if !ok {
		test.Fatalf("exp.Scope is not ast.ScopeExpression, got=%T", exp.Scope)
	}

	scope, ok = superScope.Scope.(*ast.Identifier)
	if !ok {
		test.Fatalf("superScope.Scope is not ast.Identifier, got=%T", superScope.Scope)
	}

	if !testIdentifier(test, scope, "world") {
		return
	}

	ident, ok = superScope.Value.(*ast.Identifier)
	if !ok {
		test.Fatalf("exp.Value is not ast.Identifier, got=%T", superScope.Value)
	}

	if !testIdentifier(test, ident, "compass") {
		return
	}

	ident, ok = exp.Value.(*ast.Identifier)
	if !ok {
		test.Fatalf("exp.Value is not ast.ScopeExpression, got=%T", exp.Value)
	}

	if !testIdentifier(test, ident, "south") {
		return
	}

	//

	statement, ok = program.Statements[3].(*ast.ExpressionStatement)
	if !ok {
		test.Fatalf("program.Statements[3] is not ast.ExpressionStatement, got=%T", program.Statements[3])
	}

	call, ok = statement.Expression.(*ast.CallExpression)
	if !ok {
		test.Fatalf("statement.Expression is not ast.CallExpression, got=%T", statement.Expression)
	}

	if len(call.Arguments) != 1 {
		test.Fatalf("wrong number of arguments,expected=%d, got=%T", 1, call.TokenLiteral())
	}

	if !testIntegralLiteral(test, call.Arguments[0], 1) {
		return
	}

	superScope, ok = call.Function.(*ast.ScopeExpression)
	if !ok {
		test.Fatalf("call.Function is not ast.Identifier, got=%T", call.Function)
	}

	if !testIdentifier(test, superScope.Scope, "world") {
		return
	}

	if !testIdentifier(test, superScope.Value, "Set") {
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

func testDeclarationStatement(test *testing.T, statement ast.Statement, t string, name string) bool {
	declarationStatement, ok := statement.(*ast.DeclarationStatement)
	if !ok {
		test.Errorf("statement is not *ast.DeclarationStatement type: got=%v", statement)
		return false
	}

	if statement.TokenLiteral() != "" {
		test.Errorf("statement.TokenLiteral() is not empty string: got=%v", statement.TokenLiteral())
		return false
	}

	if declarationStatement.Name.Type.Value != t {
		test.Errorf("declarationStatement.Name.Type.Value is not *%v type: got=%v", t, declarationStatement.Name.Type.Value)
		return false
	}

	if declarationStatement.Name.Type.TokenLiteral() != t {
		test.Errorf("declarationStatement.Name.Type.TokenLiteral() is not *%v type: got=%v", t, declarationStatement.Name.Type.TokenLiteral())
		return false
	}

	if declarationStatement.Name.Value != name {
		test.Errorf("declarationStatement.Name.Value is not *%v type: got=%v", name, declarationStatement.Name.Value)
		return false
	}

	if declarationStatement.Name.TokenLiteral() != name {
		test.Errorf("declarationStatement.Name.TokenLiteral is not *%v type: got=%v", name, declarationStatement.Name.TokenLiteral())
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
	case bool:
		return testBooleanLiteral(test, exp, v)
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

func testEmptyCallExpression(test *testing.T, expression ast.Expression, value string) bool {
	fun, ok := expression.(*ast.CallExpression)
	if !ok {
		test.Errorf("expression is not *ast.CallExpression, got=%T", expression)
		return false
	}

	if fun.TokenLiteral() != value {
		test.Errorf("fun.TokenLiteral() is not %v, got=%v", value, fun.TokenLiteral())
		return false
	}

	if len(fun.Arguments) != 0 {
		test.Errorf("wrong number of fun.Arguments, expected=0, got=%v", len(fun.Arguments))
		return false
	}

	ident, ok := fun.Function.(*ast.Identifier)
	if !ok {
		test.Errorf("fun.Function is not *ast.Identifier, got=%T", fun.Function)
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

func testTypedIdentifier(test *testing.T, expression ast.Expression, t string, value string) bool {
	ident, ok := expression.(*ast.TypedIdentifier)
	if !ok {
		test.Errorf("expression is not *ast.TypedIdentifier, got=%T", expression)
		return false
	}

	if !testIdentifier(test, ident.Type, t) {
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
