package tokens

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	WHITESPACE = "_"
	NEWLINE    = "NEWLINE"
	COMMA      = ","

	IDENT = "IDENT"
	BOOL  = "BOOL"
	INT   = "INT"
	DIR   = "DIR"
	USING = "USING"

	IF    = "IF"
	ELSE  = "ELSE"
	WHILE = "WHILE"

	FALSE = "FALSE"
	TRUE  = "TRUE"

	COLON  = ":"
	DCOLON = "::"
	DOLLAR = "$"

	ASSIGN = "="
	EQUAL  = "=="
	NEQUAL = "!="

	LT = "<"
	GT = ">"
	LE = "<="
	GE = ">="

	AND = "AND"
	OR  = "OR"
	NOT = "NOT"
)

var keywords = map[string]TokenType{
	"If":    IF,
	"Else":  ELSE,
	"While": WHILE,
	"True":  TRUE,
	"False": FALSE,
	"Bool":  BOOL,
	"Int":   INT,
	"Dir":   DIR,
	"Using": USING,
	"And":   AND,
	"Or":    OR,
	"Not":   NOT,
}

func LookUpIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
