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

func TestSmoke2(t *testing.T) {

	file, err := os.Open("bot_long.nil")
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
