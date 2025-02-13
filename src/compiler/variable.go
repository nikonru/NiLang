package compiler

type variable struct {
	Name name
	Addr address
	Type Type
}

type Type struct {
	Scope *scope
	Name  name
}
