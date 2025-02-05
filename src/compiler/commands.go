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

	// LOAD_TO_REG_FROM_REG [target] [source]
	LOAD_TO_REG_FROM_REG = "ld"
	LOAD_TO_REG_FROM_VAL = "ldv"
	LOAD_TO_MEM_FROM_REG = "ldr"
	LOAD_TO_REG_FROM_MEM = "ldm"

	RETURN = "ret"
)
