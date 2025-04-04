//go:build !js || !wasm

package main

import (
	"NiLang/src/common"
	"NiLang/src/compiler"
	"NiLang/src/helper"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const VERSION = "v0.1.3"
const VERSION_NAME = "Alpha"

func main() {
	stackSize := flag.Int("s", common.DefaultStackSize, "stack size in bytes")
	outputFilename := flag.String("o", "bot.tor", "output file name")
	printAST := flag.Bool("AST", false, "print abstract syntax tree in a human readable form (pseudo-code), use it for debugging the compiler")
	printVersion := flag.Bool("version", false, "print current version of the compiler")
	flag.Parse()

	if *printVersion {
		fmt.Printf("%s%s\n", VERSION, VERSION_NAME)
		return
	}

	var fileName string
	if flag.NArg() < 1 {
		log.Fatal("Expected argument with path to code to compile")
	} else {
		fileName = flag.Arg(0)
	}

	abs, err := filepath.Abs(fileName)
	if err != nil {
		log.Fatal(err)
	}
	helper.SetFilename(abs)

	file, err := os.Open(fileName)
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

	c := compiler.New(*stackSize)
	code, errors := c.Compile(input, *printAST)
	if len(errors) != 0 {
		for _, err := range errors {
			helper.PrintError(err, input)
		}
		return
	}

	output, err := os.Create(*outputFilename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := output.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	_, err = output.Write(code)
	if err != nil {
		log.Fatal(err)
	}
}
