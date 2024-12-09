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
	return nil
}
