//go:build js && wasm

package main

import (
	"NiLang/src/common"
	"NiLang/src/compiler"
	"NiLang/src/helper"
	"fmt"
	"strconv"
	"syscall/js"
)

func main() {
	js.Global().Set("compile", js.FuncOf(compile))
	js.Global().Set("getVersion", js.FuncOf(getVersion))
	select {}
}

func compile(this js.Value, args []js.Value) any {
	// returns error, string

	if len(args) < 1 {
		return js.ValueOf([]any{true, "expected source code string"})
	}

	if len(args) > 2 {
		return js.ValueOf([]any{true, fmt.Sprintf("expected no more than 2 arguments (source code and optionally stack size in bytes), got %d instead", len(args))})
	}

	source := args[0].String()
	input := []byte(source)

	stackSize := common.DefaultStackSize
	if len(args) > 1 {
		arg := args[1].String()
		val, err := strconv.Atoi(arg)
		if err != nil {
			return js.ValueOf([]any{true, fmt.Sprintf("stack size must be a number, got %v instead", arg)})
		}
		stackSize = val
	}

	c := compiler.New(stackSize)
	code, errors := c.Compile(input, false)
	if len(errors) != 0 {
		output := ""
		for _, err := range errors {
			helper.PrintError(err, input)
			output += fmt.Sprintf("%s\n", helper.FormatError(err, input))
		}

		return []any{true, string(output)}
	}

	return []any{false, string(code)}
}

func getVersion(this js.Value, args []js.Value) any {
	// returns string
	return js.ValueOf(common.VERSION + common.VERSION_NAME)
}
