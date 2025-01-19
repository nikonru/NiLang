package compiler

type name = string

type variable struct {
	Addr address
	Type name
}

const (
	Int  = "Int"
	Bool = "Bool"
)
