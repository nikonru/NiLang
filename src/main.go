package main

import (
	"NiLang/src/compiler"
	"NiLang/src/helper"
	"io"
	"log"
	"os"
)

func main() {

	args := os.Args[1:]

	outputName := "bot.tor"
	var fileName string
	if len(args) < 1 || len(args) > 2 {
		log.Fatal("Expected argument")
	} else {
		fileName = args[0]
		if len(args) > 1 {
			outputName = args[1]
		}
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

	c := compiler.New()
	code, errors := c.Compile(input)
	if len(errors) == 0 {
		for _, err := range errors {
			helper.PrintError(err, input)
		}
	}

	output, err := os.Create(outputName)
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
