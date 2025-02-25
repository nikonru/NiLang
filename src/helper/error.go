package helper

import (
	"NiLang/src/tokens"
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type Error struct {
	Line        int
	Offset      int
	Description string
}

func PrintError(error Error, input []byte) {
	fmt.Printf("%s\n", FormatError(error, input))
}

func MakeError(token tokens.Token, description string) Error {
	return Error{Line: token.Line, Offset: token.Offset, Description: description}
}

func FormatError(error Error, input []byte) (str string) {
	line := getLine(error.Line, input)
	pointer := strings.Repeat("-", len(line))

	if error.Offset < len(line) {
		index := error.Offset
		pointer = pointer[:index] + "^" + pointer[index+1:]
	}

	str = fmt.Sprintf("%s\n%s\n%s:%d:%d: %s", string(line), pointer, getFilename(), error.Line, error.Offset, error.Description)
	return str
}

var filename string

func SetFilename(name string) {
	filename = name
}

func getFilename() string {
	return filename
}

func getLine(line int, input []byte) (value []byte) {
	bytesReader := bytes.NewReader(input)
	bufReader := bufio.NewReader(bytesReader)

	for i := 0; i < line; i++ {
		value, _, _ = bufReader.ReadLine()
	}
	return value
}
