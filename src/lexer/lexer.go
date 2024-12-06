package lexer

import (
	"NiLang/src/tokens"
	"log"
)

type Lexer interface {
	NextToken() tokens.Token
}

type lexer struct {
	input []byte
	char  byte

	current int
	next    int
}

func New(input []byte) Lexer {
	l := &lexer{
		input:   input,
		current: 0,
		next:    1,
	}
	l.read()
	return l
}

func (l *lexer) NextToken() tokens.Token {
	if l.char == '#' {
		l.skipComment()
	}

	var tok tokens.Token
	switch l.char {
	case '\n', '\r':
		tok = l.newToken(tokens.NEWLINE)
		l.skipNewlines()
		return tok
	case ' ':
		tok = l.newToken(tokens.WHITESPACE)
	case ',':
		tok = l.newToken(tokens.COMMA)
	case '=':
		if l.pick() == '=' {
			tok = l.newDoubleCharacterToken(tokens.EQUAL)
		} else {
			tok = l.newToken(tokens.ASSIGN)
		}
	case '!':
		if l.pick() == '=' {
			tok = l.newDoubleCharacterToken(tokens.NEQUAL)
		} else {
			tok = l.newToken(tokens.ILLEGAL) //may be array index?
		}
	case ':':
		if l.pick() == ':' {
			tok = l.newDoubleCharacterToken(tokens.DCOLON)
		} else {
			tok = l.newToken(tokens.COLON)
		}
	case '>':
		if l.pick() == '=' {
			tok = l.newDoubleCharacterToken(tokens.GE)
		} else {
			tok = l.newToken(tokens.GT)
		}
	case '<':
		if l.pick() == '=' {
			tok = l.newDoubleCharacterToken(tokens.LE)
		} else {
			tok = l.newToken(tokens.LT)
		}
	case 0:
		tok.Literal = ""
		tok.Type = tokens.EOF
	case '\t':
		log.Fatal("Tabulation is illegal, use only spaces")
	default:
		if isDigit(l.char) {
			return l.readNumberToken()
		}
		if isLetter(l.char) {
			tok.Literal = string(l.readIdent())
			tok.Type = tokens.LookUpIdent(tok.Literal)
			return tok
		}

		tok = l.newToken(tokens.ILLEGAL)
	}

	l.read()
	return tok
}

func (l *lexer) read() {
	if l.current >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.current]
	}
	l.current = l.next
	l.next++
}

func (l *lexer) readSequence(check func(byte) bool) []byte {
	begin := l.current - 1
	end := begin
	for check(l.char) {
		l.read()
		end++
	}

	return l.input[begin:end]
}

func (l *lexer) readNumber() []byte {
	return l.readSequence(isDigit)
}

func (l *lexer) readIdent() []byte {
	return l.readSequence(isLetter)
}

func (l *lexer) readNumberToken() tokens.Token {
	return tokens.Token{
		Type:    tokens.INT,
		Literal: string(l.readNumber()),
	}
}

func (l *lexer) skipComment() {
	for l.char != '\n' && l.char != '\r' {
		l.read()
	}
	for isNewline(l.char) {
		l.read()
	}
}

func (l *lexer) skipNewlines() {
	for l.char == '\n' || l.char == '\r' {
		l.read()
	}
}

func (l *lexer) newToken(tok tokens.TokenType) tokens.Token {
	var literal string
	if isNewline(l.char) {
		literal = "newline"
	} else {
		literal = string(l.char)
	}

	return tokens.Token{
		Type:    tok,
		Literal: literal,
	}
}

func (l *lexer) newDoubleCharacterToken(tok tokens.TokenType) tokens.Token {
	first := l.char
	l.read()
	return tokens.Token{
		Type:    tok,
		Literal: string(first) + string(l.char),
	}
}

func (l *lexer) pick() byte {
	if l.next >= len(l.input) {
		return 0
	}
	return l.input[l.current]
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func isNewline(char byte) bool {
	return char == '\n' || char == '\r'
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}
