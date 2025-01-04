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

func (ds *DeclarationStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ds.TokenLiteral() + " ")
	out.WriteString(ds.Name.String())
	out.WriteString(" = ")

	if ds.Value != nil {
		out.WriteString(ds.Value.String())
	}

	return out.String()
}

type UsingStatement struct {
	Token tokens.Token
	Name  *Identifier
}

func (us *UsingStatement) statementNode()       {}
func (us *UsingStatement) TokenLiteral() string { return us.Token.Literal }

func (us *UsingStatement) String() string {
	var out bytes.Buffer

	out.WriteString(us.TokenLiteral() + " ")
	out.WriteString(us.Name.String())

	return out.String()
}

type ReturnStatement struct {
	Token tokens.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	out.WriteString(rs.Value.String())

	return out.String()
}

type ScopeStatement struct {
	Token tokens.Token
	Name  *Identifier
	Body  *BlockStatement
}

func (ss *ScopeStatement) statementNode()       {}
func (ss *ScopeStatement) TokenLiteral() string { return ss.Token.Literal }

func (ss *ScopeStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ss.TokenLiteral() + " ")
	out.WriteString(ss.Name.String() + " ")
	out.WriteString(ss.Body.String())

	return out.String()
}

type ExpressionStatement struct {
	Token      tokens.Token
	Expression Expression
}

func (ds *ExpressionStatement) statementNode()       {}
func (ds *ExpressionStatement) TokenLiteral() string { return ds.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type Identifier struct {
	Token tokens.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) statementNode()       {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegralLiteral struct {
	Token tokens.Token
	Value int64
}

func (i *IntegralLiteral) expressionNode()      {}
func (i *IntegralLiteral) statementNode()       {}
func (i *IntegralLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegralLiteral) String() string       { return i.TokenLiteral() }

type BooleanLiteral struct {
	Token tokens.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) statementNode()       {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string       { return b.TokenLiteral() }

type PrefixExpression struct {
	Token    tokens.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) statementNode()       {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }

func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    tokens.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) statementNode()       {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }

func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

type BlockStatement struct {
	Token      tokens.Token
	Statements []Statement
}

func (b *BlockStatement) statementNode()       {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }

func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range b.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type IfStatement struct {
	Token       tokens.Token
	Condition   Expression
	Consequence *BlockStatement
	Elifs       []*ElifStatement
	Alternative *BlockStatement
}

func (i *IfStatement) statementNode()       {}
func (i *IfStatement) TokenLiteral() string { return i.Token.Literal }

func (i *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	for _, elif := range i.Elifs {
		if elif == nil {
			continue
		}
		out.WriteString("elif")
		out.WriteString(elif.String())
	}

	if i.Alternative != nil {
		out.WriteString("else")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

type ElifStatement struct {
	Token       tokens.Token
	Condition   Expression
	Consequence *BlockStatement
}

func (i *ElifStatement) statementNode()       {}
func (i *ElifStatement) TokenLiteral() string { return i.Token.Literal }

func (i *ElifStatement) String() string {
	var out bytes.Buffer

	out.WriteString("elif")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	return out.String()
}

type CallExpression struct {
	Token     tokens.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) statementNode()       {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ce.Function.String())

	out.WriteString("(")
	for i, arg := range ce.Arguments {
		out.WriteString(arg.String())
		if (i + 1) != len(ce.Arguments) {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")

	return out.String()
}
