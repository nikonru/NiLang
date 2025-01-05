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
		{tokens.IDENT, "bot"},
		{tokens.NEWLINE, "newline"},
		{tokens.PIDENT, "Bool"},
		{tokens.IDENT, "hungry"},
		{tokens.ASSIGN, "="},
		{tokens.TRUE, "True"},
		{tokens.NEWLINE, "newline"},
		{tokens.WHILE, "While"},
		{tokens.IDENT, "hungry"},
		{tokens.COLON, ":"},
		{tokens.NEWLINE, "newline"},
		{tokens.INDENT, "indentation"},
		{tokens.PIDENT, "Int"},
		{tokens.IDENT, "maxEnergy"},
		{tokens.ASSIGN, "="},
		{tokens.NUMBER, "1500"},
		{tokens.NEWLINE, "newline"},
		{tokens.INDENT, "indentation"},
		{tokens.PIDENT, "ConsumeSunlight"},
		{tokens.NEWLINE, "newline"},
		{tokens.INDENT, "indentation"},
		{tokens.IF, "If"},
		{tokens.PIDENT, "GetEnergy"},
		{tokens.GT, ">"},
		{tokens.IDENT, "maxEnergy"},
		{tokens.COLON, ":"},
		{tokens.NEWLINE, "newline"},
		{tokens.INDENT, "indentation"},
		{tokens.INDENT, "indentation"},
		{tokens.PIDENT, "Fork"},
		{tokens.DOLLAR, "$"},
		{tokens.IDENT, "world"},
		{tokens.DCOLON, "::"},
		{tokens.PIDENT, "Forward"},
		{tokens.NEWLINE, "newline"},
		{tokens.INDENT, "indentation"},
		{tokens.INDENT, "indentation"},
		{tokens.IDENT, "hungry"},
		{tokens.ASSIGN, "="},
		{tokens.FALSE, "False"},
		{tokens.NEWLINE, "newline"},
		{tokens.PIDENT, "Dir"},
		{tokens.IDENT, "dir"},
		{tokens.ASSIGN, "="},
		{tokens.PIDENT, "GetDir"},
		{tokens.NEWLINE, "newline"},
		{tokens.IDENT, "bot"},
		{tokens.DCOLON, "::"},
		{tokens.PIDENT, "Move"},
		{tokens.DOLLAR, "$"},
		{tokens.PIDENT, "RotateClockwise"},
		{tokens.DOLLAR, "$"},
		{tokens.IDENT, "dir"},
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

	file, err := os.Open("bot_short.nil")
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
		Type    tokens.TokenType
		Literal string
		Line    int
		Offset  int
	}{
		{tokens.USING, "Using", 2, 0},
		{tokens.IDENT, "bot", 2, 6},
		{tokens.NEWLINE, "newline", 2, 9},
		{tokens.PIDENT, "Bool", 4, 0},
		{tokens.IDENT, "x", 4, 5},
		{tokens.ASSIGN, "=", 4, 7},
		{tokens.FALSE, "False", 4, 9},
		{tokens.NEWLINE, "newline", 4, 14},
		{tokens.PIDENT, "Ille", 5, 0},
		{tokens.ILLEGAL, ".", 5, 4},
		{tokens.IDENT, "gal", 5, 5},
		{tokens.IDENT, "y", 5, 9},
		{tokens.ASSIGN, "=", 5, 11},
		{tokens.NUMBER, "1", 5, 13},
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
