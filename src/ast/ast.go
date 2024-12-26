package ast

import (
	"NiLang/src/tokens"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type DeclarationStatement struct {
	Token tokens.Token
	Name  *Identifier
	Value Expression
}

func (ds *DeclarationStatement) statementNode()       {}
func (ds *DeclarationStatement) TokenLiteral() string { return ds.Token.Literal }

type Identifier struct {
	Token tokens.Token
	Value string
}

func (i *Identifier) statementNode()       {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type ExpressionStatement struct {
	Token      tokens.Token
	Expression Expression
}

func (ds *ExpressionStatement) statementNode()       {}
func (ds *ExpressionStatement) TokenLiteral() string { return ds.Token.Literal }

func (ds *DeclarationStatement) String() string { return "todo" }
func (es *ExpressionStatement) String() string  { return "todo" }
