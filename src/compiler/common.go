package compiler

import (
	"log"
	"slices"
)

type name = string
type address = int

const (
	BOOL_TRUE  = 1
	BOOL_FALSE = 0
)

const (
	Int  = "Int"
	Bool = "Bool"
)

const RETURN_REGISTER = AX

var BUILTIN_TYPES = []name{Int, Bool}

func builtIn(name name) Type {
	if slices.Contains(BUILTIN_TYPES, name) {
		return Type{Scope: nil, Name: name}
	}
	log.Fatalf("got unknown built in type %q", name)
	return Type{}
}
