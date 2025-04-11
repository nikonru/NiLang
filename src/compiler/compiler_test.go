package compiler_test

import (
	"NiLang/src/compiler"
	"NiLang/src/helper"
	"io"
	"log"
	"os"
	"testing"
)

const stackSize = 128

func TestSmoke1(t *testing.T) {

	file, err := os.Open("bot.nil")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	input, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	c := compiler.New(stackSize)
	_, errors := c.Compile(input, true)
	if len(errors) != 0 {
		for _, err := range errors {
			helper.PrintError(err, input)
		}
		t.Fatalf("Failed to compile code")
	}
}

func TestFailToCompileFunctionWithoutReturn(t *testing.T) {

	input := []byte(`
Fun X::Int:
    Int y = 0
Int x = X`)

	c := compiler.New(stackSize)
	_, errors := c.Compile(input, true)
	if len(errors) == 0 {
		t.Fatalf("Successfully compiled ill-formed code")
	}
}

func TestFailToCompileFunctionWithoutReturn1(t *testing.T) {

	input := []byte(`
Fun X::Int:
    Int y = 0
    If y == 1:
        Return 1
    Elif y == 0:
        Return 0

Int x = X`)

	c := compiler.New(stackSize)
	_, errors := c.Compile(input, true)
	if len(errors) == 0 {
		t.Fatalf("Successfully compiled ill-formed code")
	}
}

func TestFailToCompileFunctionWithoutReturn2(t *testing.T) {

	input := []byte(`
Fun X::Int:
    Int y = 0
    If y == 1:
        Return 1

Int x = X`)

	c := compiler.New(stackSize)
	_, errors := c.Compile(input, true)
	if len(errors) == 0 {
		t.Fatalf("Successfully compiled ill-formed code")
	}
}
