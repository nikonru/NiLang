package lexer_test

import (
	"NiLang/src/lexer"
	"NiLang/src/tokens"
	"io"
	"log"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {

	file, err := os.Open("bot.nl")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	input, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.USING, "Using"},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "bot"},
		{tokens.NEWLINE, "newline"},
		{tokens.BOOL, "Bool"},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "hungry"},
		{tokens.WHITESPACE, " "},
		{tokens.ASSIGN, "="},
		{tokens.WHITESPACE, " "},
		{tokens.TRUE, "True"},
		{tokens.NEWLINE, "newline"},
		{tokens.WHILE, "While"},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "hungry"},
		{tokens.COLON, ":"},
		{tokens.NEWLINE, "newline"},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.INT, "Int"},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "maxEnergy"},
		{tokens.WHITESPACE, " "},
		{tokens.ASSIGN, "="},
		{tokens.WHITESPACE, " "},
		{tokens.INT, "1500"},
		{tokens.NEWLINE, "newline"},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "ConsumeSunlight"},
		{tokens.NEWLINE, "newline"},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.IF, "If"},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "GetEnergy"},
		{tokens.WHITESPACE, " "},
		{tokens.GT, ">"},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "maxEnergy"},
		{tokens.COLON, ":"},
		{tokens.NEWLINE, "newline"},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "Fork"},
		{tokens.DOLLAR, "$"},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "world"},
		{tokens.DCOLON, "::"},
		{tokens.IDENT, "Forward"},
		{tokens.NEWLINE, "newline"},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "hungry"},
		{tokens.WHITESPACE, " "},
		{tokens.ASSIGN, "="},
		{tokens.WHITESPACE, " "},
		{tokens.FALSE, "False"},
		{tokens.NEWLINE, "newline"},
		{tokens.DIR, "Dir"},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "dir"},
		{tokens.WHITESPACE, " "},
		{tokens.ASSIGN, "="},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "GetDir"},
		{tokens.NEWLINE, "newline"},
		{tokens.IDENT, "bot"},
		{tokens.DCOLON, "::"},
		{tokens.IDENT, "Move"},
		{tokens.DOLLAR, "$"},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "RotateClockwise"},
		{tokens.DOLLAR, "$"},
		{tokens.WHITESPACE, " "},
		{tokens.IDENT, "dir"},
		{tokens.WHITESPACE, " "},
		{tokens.WHITESPACE, " "},
		{tokens.NEWLINE, "newline"},
	}

	Lexer := lexer.New(input)

	for i, test := range tests {
		tok := Lexer.NextToken()

		if tok.Type != test.expectedType || tok.Literal != test.expectedLiteral {
			t.Logf("tests[%d] - tok: %#v", i, tok)
			t.Fatalf("tests[%d] - token type. expected=%q, got=%q; literal. expected=%q, got=%q;", i, test.expectedType, tok.Type, test.expectedLiteral, tok.Literal)
		}
	}
}
