package compiler

type command = string

const (
	COMPARE            = "cmp"
	COMPARE_WITH_VALUE = "cmpv"

	JUMP              = "jmp"
	JUMP_IF_EQUAL     = "jme"
	JUMP_IF_NOT_EQUAL = "jne"

	LOAD     = "ld"
	LOAD_VAL = "ldv"
	LOAD_REG = "ldr"
	LOAD_MEM = "ldm"
)
