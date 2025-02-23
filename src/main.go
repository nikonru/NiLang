package main

import (
	"NiLang/src/compiler"
	"NiLang/src/helper"
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	stackSize := flag.Int("s", 128, "stack size in bytes")
	outputFilename := flag.String("o", "bot.tor", "output file name")
	printAST := flag.Bool("AST", false, "print abstract syntax tree in a human readable form (pseudo-code), use for debugging the compiler")
	flag.Parse()

	var fileName string
	if flag.NArg() < 1 {
		log.Fatal("Expected argument with path to code to compile")
	} else {
		fileName = flag.Arg(0)
	}

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
		log.Fatal("Failed to compile code")
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
