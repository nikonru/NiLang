package parser

import (
	"NiLang/src/ast"
	"NiLang/src/lexer"
	"NiLang/src/tokens"
)

type Parser struct {
	lexer *lexer.Lexer

	current tokens.Token
	next    tokens.Token
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}

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
		return false, nil
	}
}

func (p *Parser) parseDeclarationStatement() (bool, *ast.DeclarationStatement) {
	// TODO remove returning bool
	statement := &ast.DeclarationStatement{Token: p.current}

	if !p.expectNext(tokens.WHITESPACE) {
		return false, nil
	}

	p.skipWhitespaces()

	if !p.isCurrent(tokens.IDENT) {
		return false, nil
	}

	statement.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}

	if p.expectNext(tokens.WHITESPACE) {
		p.skipWhitespaces()
	}

	if !p.isCurrent(tokens.ASSIGN) {
		return false, nil
	}

	p.skipUpToNewline()

	return true, statement
}

func (p *Parser) isCurrent(t tokens.TokenType) bool {
	return t == p.current.Type
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
		return false
	}
}
