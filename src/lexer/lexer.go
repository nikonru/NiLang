package lexer

import (
	"NiLang/src/helper"
	"NiLang/src/tokens"
	"fmt"
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

	indentation bool

	line   int
	pos    int
	offset int
}

func New(input []byte) Lexer {
	l := &lexer{
		input:   input,
		current: 0,
		next:    1,

		line:   0,
		offset: 0,
	}
	l.startNewline()
	l.read()
	l.skipNewlines()
	l.offset = 0
	return l
}

func (l *lexer) NextToken() tokens.Token {
	l.pos = l.offset

	for l.char == '#' {
		flag := l.offset == 0
		l.skipComment()
		if flag {
			l.skipNewlines()
		}
	}

	var tok tokens.Token
	switch l.char {
	case '\n', '\r':
		tok = l.newToken(tokens.NEWLINE)
		l.skipNewlines()
		return tok
	case ' ':
		if l.indentation {
			offset := l.readIndent()
			return tokens.Token{Type: tokens.INDENT, Literal: "indentation", Line: l.line, Offset: offset}
		} else {
			l.read()
			return l.NextToken()
		}
	case ',':
		tok = l.newToken(tokens.COMMA)
	case '$':
		tok = l.newToken(tokens.DOLLAR)
	case '-':
		tok = l.newToken(tokens.NEGATION)
	case '=':
		if l.peek() == '=' {
			tok = l.newDoubleCharacterToken(tokens.EQUAL)
		} else {
			tok = l.newToken(tokens.ASSIGN)
		}
	case '!':
		if l.peek() == '=' {
			tok = l.newDoubleCharacterToken(tokens.NEQUAL)
		} else {
			tok = l.newToken(tokens.ILLEGAL) //may be an array index?
		}
	case ':':
		if l.peek() == ':' {
			tok = l.newDoubleCharacterToken(tokens.DCOLON)
		} else {
			tok = l.newToken(tokens.COLON)
		}
	case '>':
		if l.peek() == '=' {
			tok = l.newDoubleCharacterToken(tokens.GE)
		} else {
			tok = l.newToken(tokens.GT)
		}
	case '<':
		if l.peek() == '=' {
			tok = l.newDoubleCharacterToken(tokens.LE)
		} else {
			tok = l.newToken(tokens.LT)
		}
	case 0:
		tok = l.newToken(tokens.EOF)
		tok.Literal = ""
	case '\t':
		desc := fmt.Sprintf("tabulation is not allowed, use %d whitespaces only", tokens.INDENT_LENGTH)
		err := helper.Error{Line: l.line, Offset: l.offset, Description: desc}
		log.Fatal("\n" + helper.FormatError(err, l.input))
	default:
		if isDigit(l.char) {
			return l.readNumberToken()
		}
		if isLetter(l.char) {
			literal := string(l.readIdent())
			return l.makeToken(tokens.LookUpIdent(literal), literal)
		}

		tok = l.newToken(tokens.ILLEGAL)
	}

	l.read()
	return tok
}

func (l *lexer) startNewline() {
	l.line++
	l.offset = 0
	l.indentation = true
}

func (l *lexer) read() {
	l.offset++
	// here we need only '\n'
	if l.char == '\n' {
		l.startNewline()
	}

	if l.current >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.current]
	}
	l.current = l.next
	l.next++
	l.indentation = l.indentation && l.char == ' '
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
	check := func(char byte) bool {
		return isDigit(char) || isLetter(char)
	}
	return l.readSequence(check)
}

func (l *lexer) readIndent() int {
	counter := 0

	l.readSequence(func(char byte) bool {
		if char == ' ' {
			counter++
			return true
		}
		return false
	})

	if counter%tokens.INDENT_LENGTH != 0 {
		desc := fmt.Sprintf("expected indentation to be multiple of %d, got=%d whitespaces", tokens.INDENT_LENGTH, counter)
		err := helper.Error{Line: l.line, Offset: l.offset, Description: desc}
		log.Fatal("\n" + helper.FormatError(err, l.input))
	}

	return counter
}

func (l *lexer) readNumberToken() tokens.Token {
	return l.makeToken(tokens.NUMBER, string(l.readNumber()))
}

func (l *lexer) skipComment() {
	for !isNewline(l.char) && l.char != 0 {
		l.read()
	}
}

func (l *lexer) skipNewlines() {
	for isNewline(l.char) {
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

	return l.makeToken(tok, literal)
}

func (l *lexer) newDoubleCharacterToken(tok tokens.TokenType) tokens.Token {
	first := l.char
	l.read()
	return l.makeToken(tok, string(first)+string(l.char))
}

func (l *lexer) makeToken(tok tokens.TokenType, literal string) tokens.Token {
	return tokens.Token{Type: tok, Literal: literal, Line: l.line, Offset: l.pos}
}

func (l *lexer) peek() byte {
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
