package compiler

import "bytes"

type variable struct {
	Name name
	Addr address
	Type Type
}

type Type struct {
	Scope *scope
	Name  name
}

func (t *Type) String() string {
	var out bytes.Buffer

	scope := t.Scope
	for scope != nil {
		out.WriteString(t.Scope.name + "::")
		scope = t.Scope.GetParent()
	}
	out.WriteString(t.Name)

	return out.String()
}
