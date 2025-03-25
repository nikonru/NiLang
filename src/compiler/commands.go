package compiler

type command = string

const (
	// COMPARE [a] [b] ?a>=b ?a==b ?a<b ?a!=b
	COMPARE            = "cmp"
	COMPARE_WITH_VALUE = "cmpv"

	//JUMP [label]
	JUMP                 = "jmp"
	JUMP_IF_EQUAL        = "jme"
	JUMP_IF_NOT_EQUAL    = "jne"
	JUMP_IF_LESS_THAN    = "jml"
	JUMP_IF_GREATER_THAN = "jmg"

	JUMP_IF_LESS_EQUAL_THAN    = "jle"
	JUMP_IF_GREATER_EQUAL_THAN = "jge"

	JUMP_IF_EMPTY   = "jmf"
	JUMP_IF_FRIEND  = "jmc"
	JUMP_IF_SIBLING = "jmb"

	// LOAD_TO_REG_FROM_REG [target] [source]
	LOAD_TO_REG_FROM_REG = "ld"
	LOAD_TO_REG_FROM_VAL = "ldv"
	LOAD_TO_MEM_FROM_REG = "ldr"
	LOAD_TO_REG_FROM_MEM = "ldm"

	//CALL [label]
	CALL = "call"

	RETURN = "ret"

	MOVE = "mov"
	FACE = "rot"

	FORK             = "fork"
	SPLIT            = "split"
	BITE             = "bite"
	CONSUME_SUNLIGHT = "eatsun"
	ABSORB_MINERALS  = "absorb"
	CHECK            = "chk"
	SKIP_CYCLE       = "nop"

	NEGATE   = "neg"
	ADD      = "add"
	SUBTRACT = "sub"
	DIVIDE   = "div"
	MULTIPLY = "mul"
	MOD      = "mod"
	POWER    = "pow"
)
