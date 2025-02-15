package compiler

type function struct {
	Name      name
	Label     string
	Type      Type
	Signature []Type
}
