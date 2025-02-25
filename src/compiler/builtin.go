package compiler

import (
	"NiLang/src/ast"
	"NiLang/src/helper"
	"fmt"
	"log"
)

func (c *Compiler) initBuiltin(globalScope *scope) {
	builtins := []struct {
		name              string
		numberOfArguments int
	}{
		{"Fork", 1},
		{"Split", 1},
		{"Bite", 1},
		{"ConsumeLight", 0},
		{"AbsorbMinerals", 0},
		{"IsEmpty", 1},
		{"IsSibling", 1},
		{"IsFriend", 1},
		{"GetLuminosity", 1},
		{"GetMineralization", 1},
		{"Sleep", 0},
		{"Move", 1},
		{"Face", 1},
	}
	bot := newScope("bot")

	for _, builtin := range builtins {
		bot.functions[builtin.name] = function{
			Name:      builtin.name,
			Label:     "",
			Type:      VOID, //we shouldn't check this at all
			Arguments: make([]variable, builtin.numberOfArguments),
			IsBuiltin: true}
	}

	ok := globalScope.AddScope(bot)
	if !ok {
		log.Fatalf("failed to initialize builtin variables")
	}

	dir := newScope(helper.FirstToLowerCase(Dir))

	var directions = [DIR_END]string{"_", "front", "frontRight", "right", "backRight", "back", "backLeft", "left", "frontLeft"}

	for direction := DIR_BEGIN + 1; direction < DIR_END; direction++ {
		addr := c.purchaseMemoryAddress()
		ok := dir.AddVariable(directions[direction], addr, builtIn(Dir))
		if !ok {
			log.Fatalf("failed to initialize builtin variables")
		}
		c.emit(LOAD_TO_REG_FROM_VAL, AX, direction)
		c.emit(LOAD_TO_MEM_FROM_REG, addr, AX)
	}

	ok = globalScope.AddScope(dir)
	if !ok {
		log.Fatalf("failed to initialize builtin variables")
	}
}

func (c *Compiler) compileBuiltin(expression *ast.CallExpression, name name) (Type, register) {
	emitComparison := func(condition command) (Type, register) {
		True := c.getUniqueLabel()
		end := c.getUniqueLabel()

		c.emit(condition, True)
		c.emit(LOAD_TO_REG_FROM_VAL, AX, BOOL_FALSE)
		c.emit(JUMP, end)

		c.emitLabel(True)
		c.emit(LOAD_TO_REG_FROM_VAL, AX, BOOL_TRUE)
		c.emitLabel(end)

		return builtIn(Bool), AX
	}

	direction := func() register {
		numberOfArguments := 1
		if len(expression.Arguments) != numberOfArguments {
			err := helper.MakeError(expression.Token,
				fmt.Sprintf("unexpected number of arguments expected=%d, got=%d", numberOfArguments, len(expression.Arguments)))
			c.addError(err)
			return ""
		}
		t, register := c.compileExpression(expression.Arguments[0])

		if t != builtIn(Dir) {
			err := helper.MakeError(expression.Token,
				fmt.Sprintf("unexpected type of an argument expected %q, got %q", Dir, t.String()))
			c.addError(err)
		}
		return register
	}

	switch name {
	case "Fork":
		c.compileFunctionWithDirectionArgument(FORK, direction())
		return VOID, ""
	case "Split":
		c.compileFunctionWithDirectionArgument(SPLIT, direction())
		return VOID, ""
	case "Bite":
		c.compileFunctionWithDirectionArgument(BITE, direction())
		return VOID, ""
	case "ConsumeSunlight":
		c.emit(CONSUME_SUNLIGHT)
		return VOID, ""
	case "AbsorbMinerals":
		c.emit(ABSORB_MINERALS)
		return VOID, ""
	case "IsEmpty":
		c.compileFunctionWithDirectionArgument(CHECK, direction())
		return emitComparison(JUMP_IF_EMPTY)
	case "IsSibling":
		c.compileFunctionWithDirectionArgument(CHECK, direction())
		return emitComparison(JUMP_IF_SIBLING)
	case "IsFriend":
		c.compileFunctionWithDirectionArgument(CHECK, direction())
		return emitComparison(JUMP_IF_FRIEND)
	case "GetLuminosity":
		c.compileFunctionWithDirectionArgument(CHECK, direction())
		c.emit(LOAD_TO_REG_FROM_REG, AX, SD)
		return builtIn(Int), AX
	case "GetMineralization":
		c.compileFunctionWithDirectionArgument(CHECK, direction())
		c.emit(LOAD_TO_REG_FROM_REG, AX, MD)
		return builtIn(Int), AX
	case "Sleep":
		c.emit(SKIP_CYCLE)
		return VOID, ""
	case "Move":
		c.compileFunctionWithDirectionArgument(MOVE, direction())
		return VOID, ""
	case "Face":
		c.compileFunctionWithDirectionArgument(FACE, direction())
		return VOID, ""
	default:
		log.Fatalf("builtin function %q is not handled", name)
		return VOID, ""
	}
}

func (c *Compiler) compileFunctionWithDirectionArgument(command command, register register) {
	var labels [DIR_END]string
	for dir := DIR_BEGIN + 1; dir < DIR_END; dir++ {
		c.emit(COMPARE_WITH_VALUE, register, dir)
		label := c.getUniqueLabel()
		c.emit(JUMP_IF_EQUAL, label)

		labels[dir] = label
	}

	end := c.getUniqueLabel()

	for dir := DIR_BEGIN + 1; dir < DIR_END; dir++ {
		c.emitLabel(labels[dir])

		var direction string
		switch dir {
		case FRONT:
			direction = "front"
		case FRONT_RIGHT:
			direction = "frontright"
		case RIGHT:
			direction = "right"
		case BACK_RIGHT:
			direction = "backright"
		case BACK:
			direction = "back"
		case BACK_LEFT:
			direction = "backleft"
		case LEFT:
			direction = "left"
		case FRONT_LEFT:
			direction = "frontleft"
		}

		if command == SPLIT || command == FORK {
			c.emit(command, direction, BEGIN_LABEL)
		} else {
			c.emit(command, direction)
		}
		c.emit(JUMP, end)
	}

	c.emitLabel(end)
}
