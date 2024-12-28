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

	file, err := os.Open("bot.nil")
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
		{tokens.NUMBER, "1500"},
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

func TestLexerLines(t *testing.T) {

	input := []byte(`Using bot

Bool x = False
Ille.gal y = 1`)

	tests := []struct {
		Type    tokens.TokenType
		Literal string
		Line    int
		Offset  int
	}{
		{tokens.USING, "Using", 1, 0},
		{tokens.WHITESPACE, " ", 1, 5},
		{tokens.IDENT, "bot", 1, 6},
		{tokens.NEWLINE, "newline", 1, 9},
		{tokens.BOOL, "Bool", 3, 0},
		{tokens.WHITESPACE, " ", 3, 4},
		{tokens.IDENT, "x", 3, 5},
		{tokens.WHITESPACE, " ", 3, 6},
		{tokens.ASSIGN, "=", 3, 7},
		{tokens.WHITESPACE, " ", 3, 8},
		{tokens.FALSE, "False", 3, 9},
		{tokens.NEWLINE, "newline", 3, 14},
		{tokens.IDENT, "Ille", 4, 0},
		{tokens.ILLEGAL, ".", 4, 4},
		{tokens.IDENT, "gal", 4, 5},
		{tokens.WHITESPACE, " ", 4, 8},
		{tokens.IDENT, "y", 4, 9},
		{tokens.WHITESPACE, " ", 4, 10},
		{tokens.ASSIGN, "=", 4, 11},
		{tokens.WHITESPACE, " ", 4, 12},
		{tokens.NUMBER, "1", 4, 13},
	}

	Lexer := lexer.New(input)

	for i, test := range tests {
		tok := Lexer.NextToken()

		if tok.Type != test.Type || tok.Literal != test.Literal || tok.Line != test.Line || tok.Offset != test.Offset {
			t.Logf("tests[%d] - tok: %#v", i, tok)
			t.Fatalf("tests[%d] - token type. expected=%q, got=%q; literal. expected=%q, got=%q; line. expected=%d, got=%d; offset. expected=%d, got=%d;",
				i, test.Type, tok.Type, test.Literal, tok.Literal, test.Line, tok.Line, test.Offset, tok.Offset)
		}
	}
}
