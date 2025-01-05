package tokens

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Offset  int
}

const INDENT_LENGTH = 4

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	INDENT  = "INDENTATION"
	NEWLINE = "NEWLINE"
	COMMA   = ","

	IDENT  = "IDENT"
	NUMBER = "NUMBER"

	USING  = "USING"
	IF     = "IF"
	ELSE   = "ELSE"
	ELIF   = "ELIF"
	WHILE  = "WHILE"
	RETURN = "RETURN"
	SCOPE  = "SCOPE"
	ALIAS  = "ALIAS"
	FUN    = "FUN"

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
	"If":     IF,
	"Else":   ELSE,
	"Elif":   ELIF,
	"While":  WHILE,
	"True":   TRUE,
	"False":  FALSE,
	"Using":  USING,
	"And":    AND,
	"Or":     OR,
	"Not":    NOT,
	"Return": RETURN,
	"Scope":  SCOPE,
	"Alias":  ALIAS,
	"Fun":    FUN,
}

func LookUpIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
