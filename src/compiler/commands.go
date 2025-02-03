package compiler

type command = string

const (
	COMPARE            = "cmp"
	COMPARE_WITH_VALUE = "cmpv"

	JUMP                 = "jmp"
	JUMP_IF_EQUAL        = "jme"
	JUMP_IF_NOT_EQUAL    = "jne"
	JUMP_IF_LESS_THAN    = "jml"
	JUMP_IF_GREATER_THAN = "jmg"

	JUMP_IF_LESS_EQUAL_THAN    = "jle"
	JUMP_IF_GREATER_EQUAL_THAN = "jge"

	LOAD     = "ld"
	LOAD_VAL = "ldv"
	LOAD_REG = "ldr"
	LOAD_MEM = "ldm"
)
