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
	PIDENT = "PIDENT" //prime identifier starts with uppercase letter
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

	AND = "And"
	OR  = "Or"
	NOT = "Not"
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
	if len(ident) == 0 {
		panic("empty identity")
	}

	if tok, ok := keywords[ident]; ok {
		return tok
	}

	if isUppercase(ident[0]) {
		return PIDENT
	}
	return IDENT
}

func isUppercase(char byte) bool {
	return 'A' <= char && char <= 'Z'
}

func GetIdentLevel(t Token) int {
	return t.Offset / INDENT_LENGTH
}
