package compiler

type Compiler struct {
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(input []byte) ([]byte, error) {
	return input, nil
}
