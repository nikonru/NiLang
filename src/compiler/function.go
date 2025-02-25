package compiler

type function struct {
	Name      name
	Label     string
	Type      Type
	Arguments []variable

	IsBuiltin bool
}
