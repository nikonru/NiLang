package parser

import (
	"NiLang/src/ast"
	"NiLang/src/helper"
	"NiLang/src/lexer"
	"NiLang/src/tokens"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	LOGIC       // And, Or
	EQUALS      // ==
	LESSGREATER // > or <
	PREFIX      // Not
	CALL        // func$ or func
	SCOPE       // ::
)

var precedence = map[tokens.TokenType]int{
	tokens.DCOLON: SCOPE,
	tokens.DOLLAR: CALL,
	tokens.EQUAL:  EQUALS,
	tokens.NEQUAL: EQUALS,
	tokens.LT:     LESSGREATER,
	tokens.LE:     LESSGREATER,
	tokens.GT:     LESSGREATER,
	tokens.GE:     LESSGREATER,
	tokens.OR:     LOGIC,
	tokens.AND:    LOGIC,
}

type errors = []helper.Error
type prefixParseFns = func() ast.Expression
type infixParseFns = func(ast.Expression) ast.Expression

type Parser struct {
	lexer  *lexer.Lexer
	errors errors

	current tokens.Token
	next    tokens.Token
	level   int

	prefixParseFns map[tokens.TokenType]prefixParseFns
	infixParseFns  map[tokens.TokenType]infixParseFns

	pleaseDontSkipToken bool
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, level: 0}

	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFns)
	p.registerPrefix(tokens.IDENT, p.parseIdentifier)
	p.registerPrefix(tokens.PIDENT, p.parseCallExpressionPrefix)
	p.registerPrefix(tokens.NUMBER, p.parseIntegralLiteral)
	p.registerPrefix(tokens.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.NOT, p.parsePrefixExpression)

	p.infixParseFns = make(map[tokens.TokenType]infixParseFns)
	p.registerInfix(tokens.LT, p.parseInfixExpression)
	p.registerInfix(tokens.LE, p.parseInfixExpression)
	p.registerInfix(tokens.GT, p.parseInfixExpression)
	p.registerInfix(tokens.GE, p.parseInfixExpression)
	p.registerInfix(tokens.NEQUAL, p.parseInfixExpression)
	p.registerInfix(tokens.EQUAL, p.parseInfixExpression)
	p.registerInfix(tokens.AND, p.parseInfixExpression)
	p.registerInfix(tokens.OR, p.parseInfixExpression)
	p.registerInfix(tokens.DOLLAR, p.parseCallExpression)
	p.registerInfix(tokens.DCOLON, p.parseScopeExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	if p.isCurrent(tokens.NEWLINE) {
		p.level = 0
	}

	p.current = p.next
	p.next = (*p.lexer).NextToken()

	if p.isCurrent(tokens.INDENT) {
		p.level = tokens.GetIdentLevel(p.current)
		p.nextToken()
	}
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.isCurrent(tokens.EOF) {
		ok, statement := p.parseStatement()
		if ok {
			program.Statements = append(program.Statements, statement)
		}

		if p.pleaseDontSkipToken {

			p.pleaseDontSkipToken = false
			continue
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() (bool, ast.Statement) {
	switch p.current.Type {
	case tokens.USING:
		return p.parseUsingStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	case tokens.SCOPE:
		return p.parseScopeStatement()
	case tokens.WHILE:
		return p.parseWhileStatement()
	case tokens.IF:
		return p.parseIfStatement()
	case tokens.ALIAS:
		return p.parseAliasStatement()
	case tokens.FUN:
		return p.parseFunctionStatement()
	case tokens.COLON, tokens.EOF, tokens.INDENT, tokens.NEWLINE, tokens.ELIF:
		err := helper.MakeError(p.current, fmt.Sprintf("attempt to parse invalid token %s", p.current.Type))
		p.addError(err)
		return false, nil
	default:
		if p.isCurrent(tokens.IDENT) && p.isNext(tokens.ASSIGN) {
			return p.parseAssignmentStatement()
		}

		if p.isCurrent(tokens.PIDENT) && p.isNext(tokens.IDENT) {
			return p.parseDeclarationStatement(true)
		}

		ok, expression := p.parseExpressionStatement()

		if _type, ok := expression.Expression.(*ast.CallExpression); ok {
			scope, isScopeExpression := _type.Function.(*ast.ScopeExpression)

			isOnTheSameLine := p.current.Line == p.next.Line
			if p.isNext(tokens.IDENT) && isOnTheSameLine && len(_type.Arguments) == 0 && isScopeExpression {
				ok, ds := p.parseDeclarationStatement(false)
				if ok {
					ds.Var.Type = scope
				}
				return ok, ds
			}
		}

		return ok, expression
	}
}

func (p *Parser) parseDeclarationStatement(parseType bool) (bool, *ast.DeclarationStatement) {
	var t ast.Expression
	if parseType {
		t = p.parseType()
	}

	if !p.expectNext(tokens.IDENT) {
		return false, nil
	}

	ident := ast.Variable{Token: p.current, Type: t, Name: p.current.Literal}

	if !p.expectNext(tokens.ASSIGN) {
		return false, nil
	}

	statement := &ast.DeclarationStatement{Var: ident, Value: nil}
	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)

	return true, statement
}

func (p *Parser) parseType() ast.Expression {
	if p.isCurrent(tokens.PIDENT) {
		return &ast.Identifier{Token: p.current, Value: p.current.Literal}
	} else if p.isCurrent(tokens.IDENT) {
		ok, exp := p.parseExpressionStatement()
		if !ok {
			error := helper.MakeError(p.current, "couldn't parse type expression")
			p.addError(error)
		}

		switch t := exp.Expression.(type) {
		case *ast.ScopeExpression:
			return t
		case *ast.CallExpression:
			scope, isScopeExpression := t.Function.(*ast.ScopeExpression)
			if isScopeExpression {
				return scope
			}
		default:
			error := helper.MakeError(p.current, fmt.Sprintf("expected scope expression, got=%T", exp.Expression))
			p.addError(error)
		}
	}

	error := helper.MakeError(p.current, fmt.Sprintf("type starts with the wrong token expected %q or %q, got=%q",
		tokens.PIDENT, tokens.IDENT, p.current))
	p.addError(error)

	return nil
}

func (p *Parser) parseUsingStatement() (bool, *ast.UsingStatement) {
	statement := &ast.UsingStatement{Token: p.current}

	if !p.expectNext(tokens.IDENT) {
		return false, nil
	}

	statement.Name = p.parseExpression(LOWEST)

	return true, statement
}

func (p *Parser) parseReturnStatement() (bool, *ast.ReturnStatement) {
	statement := &ast.ReturnStatement{Token: p.current}

	p.nextToken()

	if p.isCurrent(tokens.NEWLINE) || p.isCurrent(tokens.EOF) {
		statement.Value = nil
		return true, statement
	}

	statement.Value = p.parseExpression(LOWEST)

	return true, statement
}

func (p *Parser) parseScopeStatement() (bool, *ast.ScopeStatement) {
	statement := &ast.ScopeStatement{Token: p.current}

	if !p.expectNext(tokens.IDENT) {
		return false, nil
	}

	statement.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}

	if !p.gotoBlockStatement() {
		return false, nil
	}

	statement.Body = p.parseBlockStatement()

	return true, statement
}

func (p *Parser) parseWhileStatement() (bool, *ast.WhileStatement) {
	statement := &ast.WhileStatement{Token: p.current}

	p.nextToken()
	statement.Condition = p.parseExpression(LOWEST)

	if !p.gotoBlockStatement() {
		return false, nil
	}

	statement.Body = p.parseBlockStatement()

	return true, statement
}

func (p *Parser) parseAliasStatement() (bool, *ast.AliasStatement) {
	statement := &ast.AliasStatement{Token: p.current}

	if !p.expectNext(tokens.PIDENT) {
		return false, nil
	}

	statement.Var = ast.Variable{Token: p.current, Name: p.current.Literal}

	if !p.expectNext(tokens.DCOLON) {
		return false, nil
	}

	p.nextToken()

	statement.Var.Type = p.parseType()

	if !p.gotoBlockStatement() {
		return false, nil
	}

	statement.Values = p.parseAliasValues(statement.Var.Type)
	if statement.Values == nil {
		return false, statement
	}

	return true, statement
}

func (p *Parser) parseFunctionStatement() (bool, *ast.FunctionStatement) {
	statement := &ast.FunctionStatement{Token: p.current}

	if !p.expectNext(tokens.PIDENT) {
		return false, nil
	}

	name := ast.Variable{Token: p.current, Name: p.current.Literal}

	if p.isNext(tokens.DCOLON) {
		p.nextToken()
		if p.isNext(tokens.PIDENT) || p.isNext(tokens.IDENT) {
			p.nextToken()
			name.Type = p.parseType()
		}
	} else {
		name.Type = nil
	}

	statement.Var = name

	if p.isNext(tokens.DOLLAR) {
		p.nextToken()
		statement.Parameters = p.parseFunctionParameters()
	} else {
		statement.Parameters = nil
	}

	if !p.gotoBlockStatement() {
		return false, nil
	}

	statement.Body = p.parseBlockStatement()

	return true, statement
}

func (p *Parser) parseAssignmentStatement() (bool, *ast.AssignmentStatement) {
	name := &ast.Identifier{Token: p.current, Value: p.current.Literal}
	if !p.expectNext(tokens.ASSIGN) {
		return false, nil
	}
	statement := &ast.AssignmentStatement{Name: name}

	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)

	return true, statement
}

func (p *Parser) parseExpressionStatement() (bool, *ast.ExpressionStatement) {
	statement := &ast.ExpressionStatement{Token: p.current}
	statement.Expression = p.parseExpression(LOWEST)
	return true, statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixParseFns[p.current.Type]
	if !ok {
		error := helper.MakeError(p.current, fmt.Sprintf("no prefix parse function for %s found", p.current.Type))
		p.addError(error)
		return nil
	}
	leftExpression := prefix()

	for !p.isNext(tokens.NEWLINE) && precedence < p.nextPrecendence() {
		p.nextToken()

		infix, ok := p.infixParseFns[p.current.Type]
		if !ok {
			return leftExpression
		}

		leftExpression = infix(leftExpression)
	}

	if p.isNext(tokens.NEWLINE) {
		p.nextToken()
	}

	return leftExpression
}

func (p *Parser) isCurrent(t tokens.TokenType) bool {
	return t == p.current.Type
}

func (p *Parser) isNext(t tokens.TokenType) bool {
	return t == p.next.Type
}

func (p *Parser) IsNextLevel() int {
	if p.isNext(tokens.INDENT) {
		return tokens.GetIdentLevel(p.next)
	}
	return 0
}

func (p *Parser) expectNext(t tokens.TokenType) bool {
	if p.isNext(t) {
		p.nextToken()
		return true
	} else {
		p.nextError(t)
		return false
	}
}

func (p *Parser) Errors() errors { return p.errors }

func (p *Parser) error(token tokens.TokenType, expected tokens.Token, name string) {
	desc := fmt.Sprintf("expected %s token to be %s, got %s instead", name, token, expected.Type)
	error := helper.MakeError(expected, desc)
	p.addError(error)
}

func (p *Parser) addError(error helper.Error) {
	p.errors = append(p.errors, error)
}

func (p *Parser) nextError(token tokens.TokenType) {
	p.error(token, p.next, "next")
}

func (p *Parser) registerPrefix(tokenType tokens.TokenType, fn prefixParseFns) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType tokens.TokenType, fn infixParseFns) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseIntegralLiteral() ast.Expression {
	lit := &ast.IntegralLiteral{Token: p.current}

	value, err := strconv.ParseInt(p.current.Literal, 0, 64)
	if err != nil {
		error := helper.MakeError(p.current, fmt.Sprintf("could not parse %q as integer", p.current.Literal))
		p.addError(error)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	lit := &ast.BooleanLiteral{Token: p.current}

	tokenType := tokens.LookUpIdent(p.current.Literal)
	if tokenType == tokens.TRUE {
		lit.Value = true
	} else if tokenType == tokens.FALSE {
		lit.Value = false
	} else {
		error := helper.MakeError(p.current, fmt.Sprintf("could not parse %q as boolean", p.current.Literal))
		p.addError(error)
		return nil
	}

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.current, Operator: p.current.Literal}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.current,
		Operator: p.current.Literal,
		Left:     left,
	}
	precedence := p.currentPrecendence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) gotoBlockStatement() bool {
	if !p.expectNext(tokens.COLON) {
		return false
	}

	if !p.expectNext(tokens.NEWLINE) {
		return false
	}
	return true
}

func (p *Parser) parseIfStatement() (bool, ast.Statement) {
	statement := &ast.IfStatement{Token: p.current}
	statement.Elifs = make([]*ast.ElifStatement, 0)

	level := p.level

	p.nextToken()
	statement.Condition = p.parseExpression(LOWEST)

	if !p.gotoBlockStatement() {
		return false, nil
	}

	statement.Consequence = p.parseBlockStatement()
	p.nextToken()

	for p.isCurrent(tokens.ELIF) {
		p.nextToken()

		exp := &ast.ElifStatement{Token: p.current}
		exp.Condition = p.parseExpression(LOWEST)
		if !p.gotoBlockStatement() {
			exp = nil
		}

		if exp != nil {
			exp.Consequence = p.parseBlockStatement()
		}

		statement.Elifs = append(statement.Elifs, exp)
		p.nextToken()
	}

	if p.isCurrent(tokens.ELSE) {
		if !p.gotoBlockStatement() {
			return false, nil
		}

		statement.Alternative = p.parseBlockStatement()
	}

	if p.isCurrent(tokens.NEWLINE) {
		p.nextToken()
	}
	p.pleaseDontSkipToken = true

	if p.level > level && !p.isCurrent(tokens.EOF) {
		err := helper.MakeError(p.current, "unexpected indentation after if statement")
		p.addError(err)
	}

	return true, statement
}

func (p *Parser) gotoNextLine() bool {
	if p.isCurrent(tokens.NEWLINE) {
		if p.IsNextLevel() == p.level {
			p.nextToken()
			return true
		}
	}
	return false
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}

	level := p.level
	p.nextToken()
	level += 1

	if level != p.level {
		desc := fmt.Sprintf("expected one level of indentation after expression, got %d instead", p.level)
		error := helper.MakeError(p.current, desc)
		p.addError(error)
		return nil
	}

	isInBlock := func() bool {
		return level == p.level && !p.isCurrent(tokens.EOF)
	}

	for isInBlock() {
		ok, statement := p.parseStatement()
		if ok {
			block.Statements = append(block.Statements, statement)
		}

		if !isInBlock() {
			break
		}

		if p.pleaseDontSkipToken {
			p.pleaseDontSkipToken = false
			continue
		}

		if !p.gotoNextLine() {
			break
		}
	}
	return block
}

func (p *Parser) parseAliasValues(t ast.Expression) []*ast.DeclarationStatement {
	statements := []*ast.DeclarationStatement{}

	level := p.level
	p.nextToken()
	level += 1

	if level != p.level {
		desc := fmt.Sprintf("expected one level of indentation after expression, got %d instead", p.level)
		error := helper.MakeError(p.current, desc)
		p.addError(error)
		return nil
	}

	for level == p.level && !p.isCurrent(tokens.EOF) {
		if !p.isCurrent(tokens.IDENT) {
			p.error(tokens.IDENT, p.current, "current")
			return nil
		}

		name := ast.Variable{Token: p.current, Type: t, Name: p.current.Literal}
		declaration := &ast.DeclarationStatement{Var: name}

		if !p.expectNext(tokens.ASSIGN) {
			return nil
		}

		p.nextToken()
		declaration.Value = p.parseExpression(LOWEST)

		statements = append(statements, declaration)

		if !p.gotoNextLine() {
			break
		}
	}

	return statements
}

func (p *Parser) parseCallExpressionPrefix() ast.Expression {
	fun := p.parseIdentifier()
	if p.isNext(tokens.DOLLAR) {
		return fun
	}
	exp := &ast.CallExpression{Token: p.current, Function: fun}
	exp.Arguments = nil
	return exp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.current, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.isNext(tokens.NEWLINE) {
		// maybe we never visit this if
		p.nextToken()
		return nil
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.isNext(tokens.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	return args
}

func (p *Parser) parseFunctionParameters() []ast.Variable {
	parameters := []ast.Variable{}

	for !p.isNext(tokens.COLON) {
		if len(parameters) != 0 {
			if !p.expectNext(tokens.COMMA) {
				return nil
			}
		}

		if !p.expectNext(tokens.IDENT) {
			return nil
		}
		parameter := ast.Variable{Token: p.current, Name: p.current.Literal}

		p.nextToken()
		parameter.Type = p.parseType()

		parameters = append(parameters, parameter)
	}

	return parameters
}

func (p *Parser) parseScopeExpression(scope ast.Expression) ast.Expression {
	exp := &ast.ScopeExpression{Token: p.current, Scope: scope}
	if p.isNext(tokens.IDENT) || p.isNext(tokens.PIDENT) {
		p.nextToken()
	} else {
		return nil
	}

	exp.Value = p.parseIdentifier().(*ast.Identifier)

	if p.isCurrent(tokens.PIDENT) {
		if p.isNext(tokens.DOLLAR) {
			p.nextToken()
			return p.parseCallExpression(exp)
		}
		exp := &ast.CallExpression{Token: p.current, Function: exp}
		exp.Arguments = nil
		return exp
	}

	return exp
}

func (p *Parser) nextPrecendence() int {
	if pred, ok := precedence[p.next.Type]; ok {
		return pred
	}

	return LOWEST
}

func (p *Parser) currentPrecendence() int {
	if pred, ok := precedence[p.current.Type]; ok {
		return pred
	}

	return LOWEST
}
