package compiler

type name = string

type variable struct {
	Name name
	Addr address
	Type name
}

const (
	BOOL_TRUE  = 1
	BOOL_FALSE = 0
)

const (
	Int  = "Int"
	Bool = "Bool"
)
