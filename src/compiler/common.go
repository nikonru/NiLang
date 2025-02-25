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
	Dir  = "Dir"
)

const (
	DIR_BEGIN int = iota
	FRONT
	FRONT_RIGHT
	RIGHT
	BACK_RIGHT
	BACK
	BACK_LEFT
	LEFT
	FRONT_LEFT
	DIR_END
)

const RETURN_REGISTER = AX

const BEGIN_LABEL = "BEGIN"

var BUILTIN_TYPES = []name{Int, Bool, Dir}

func builtIn(name name) Type {
	if slices.Contains(BUILTIN_TYPES, name) {
		return Type{Scope: nil, Name: name}
	}
	log.Fatalf("got unknown built in type %q", name)
	return Type{}
}
