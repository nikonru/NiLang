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
)

var precedence = map[tokens.TokenType]int{
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

	prefixParseFns map[tokens.TokenType]prefixParseFns
	infixParseFns  map[tokens.TokenType]infixParseFns
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}

	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFns)
	p.registerPrefix(tokens.IDENT, p.parseIdentifier)
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

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.current = p.next
	p.next = (*p.lexer).NextToken()
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.current.Type != tokens.EOF {
		ok, statement := p.parseStatement()
		if ok {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() (bool, ast.Statement) {
	switch p.current.Type {
	case tokens.BOOL, tokens.DIR, tokens.INT:
		return p.parseDeclarationStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDeclarationStatement() (bool, *ast.DeclarationStatement) {
	statement := &ast.DeclarationStatement{Token: p.current}

	if !p.expectNext(tokens.IDENT) {
		return false, nil
	}

	statement.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}

	if !p.expectNext(tokens.ASSIGN) {
		return false, nil
	}

	p.skipUpToNewline()

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

	for p.isNext(tokens.INDENT) {
		p.nextToken()
	}

	for !p.isNext(tokens.NEWLINE) && precedence < p.nextPrecendence() {
		p.nextToken()

		infix, ok := p.infixParseFns[p.current.Type]
		if !ok {
			return leftExpression
		}

		leftExpression = infix(leftExpression)
	}

	for p.isNext(tokens.INDENT) {
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

func (p *Parser) skipUpToNewline() {
	for !(p.isCurrent(tokens.NEWLINE) || p.isCurrent(tokens.EOF)) {
		p.nextToken()
	}
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
