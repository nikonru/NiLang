package compiler

import (
	"bytes"
)

type variable struct {
	Name name
	Addr address
	Type Type
}

type Type struct {
	Scope *scope
	Name  name
}

var VOID = Type{Scope: nil, Name: ""}

func (t *Type) String() string {
	var out bytes.Buffer

	names := make([]string, 0)
	scope := t.Scope
	for scope != nil {
		if scope.name != "" {
			names = append(names, scope.name)
		}
		scope = scope.GetParent()
	}

	for i := len(names) - 1; i >= 0; i -= 1 {
		out.WriteString(names[i] + "::")
	}

	out.WriteString(t.Name)
	return out.String()
}
