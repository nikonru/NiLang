package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	LET   = "LET"

	IF    = "IF"
	ELSE  = "ELSE"
	WHILE = "WHILE"

	FALSE = "FALSE"
	TRUE  = "TRUE"
)

var keywords = map[string]TokenType{
	"if":    IF,
	"else":  ELSE,
	"True":  TRUE,
	"False": FALSE,
}

func LookUpIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
