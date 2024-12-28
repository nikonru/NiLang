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
	EQUALS      // ==
	LESSGREATER // > or <
	PREFIX      // Not
	CALL        // func$ or func
)

type errors = []helper.Error
type prefixParseFns = func() ast.Expression

type Parser struct {
	lexer  *lexer.Lexer
	errors errors

	current tokens.Token
	next    tokens.Token

	prefixParseFns map[tokens.TokenType]prefixParseFns
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}

	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFns)
	p.registerPrefix(tokens.IDENT, p.parseIdentifier)
	p.registerPrefix(tokens.NUMBER, p.parseIntegralLiteral)
	p.registerPrefix(tokens.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.NOT, p.parsePrefixExpression)

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

	if !p.expectNext(tokens.WHITESPACE) {
		return false, nil
	}

	p.skipWhitespaces()

	if !p.expectCurrent(tokens.IDENT) {
		return false, nil
	}

	statement.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}

	if p.expectNext(tokens.WHITESPACE) {
		p.skipWhitespaces()
	}

	if !p.expectCurrent(tokens.ASSIGN) {
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
	prefix := p.prefixParseFns[p.current.Type]
	if prefix == nil {
		error := helper.MakeError(p.current, fmt.Sprintf("no prefix parse function for %s found", p.current.Type))
		p.addError(error)
		return nil
	}
	leftExpression := prefix()

	return leftExpression
}

func (p *Parser) isCurrent(t tokens.TokenType) bool {
	return t == p.current.Type
}

func (p *Parser) expectCurrent(t tokens.TokenType) bool {
	if p.isCurrent(t) {
		return true
	} else {
		p.currentError(t)
		return false
	}
}

func (p *Parser) isNext(t tokens.TokenType) bool {
	return t == p.next.Type
}

func (p *Parser) skipWhitespaces() {
	for p.isCurrent(tokens.WHITESPACE) {
		p.nextToken()
	}
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

func (p *Parser) currentError(token tokens.TokenType) {
	p.error(token, p.current, "current")
}

func (p *Parser) registerPrefix(tokenType tokens.TokenType, fn prefixParseFns) {
	p.prefixParseFns[tokenType] = fn
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

	if p.expectNext(tokens.WHITESPACE) {
		p.skipWhitespaces()
	}

	expression.Right = p.parseExpression(PREFIX)

	return expression
}
