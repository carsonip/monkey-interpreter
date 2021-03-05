package ast

import (
	"fmt"
	"github.com/carsonip/monkey-interpreter/token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
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

type Identifier struct {
	Token token.Token
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) expression() {}

type LetStatement struct {
	Token token.Token
	Name *Identifier
	Value Expression
}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

func (l *LetStatement) statement() {}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) statement() {}

type IfStatement struct {
	Token token.Token
	Condition Expression
	Then []Node
	Else []Node
}

func (s *IfStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *IfStatement) statement() {}

type NumberLiteral struct {
	Token token.Token
	Value int
}

func (n *NumberLiteral) TokenLiteral() string {
	return n.Token.Literal
}

func (n *NumberLiteral) expression() {}

type InfixExpression struct {
	Token token.Token
	Left Expression
	Right Expression
}

func (i *InfixExpression) TokenLiteral() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.TokenLiteral(), i.Token.Literal, i.Right.TokenLiteral())
}

func (i *InfixExpression) expression() {}

type PrefixExpression struct {
	Token token.Token
	Right Expression
}

func (p *PrefixExpression) TokenLiteral() string {
	return fmt.Sprintf("(%s%s)", p.Token.Literal, p.Right.TokenLiteral())
}

func (p *PrefixExpression) expression() {}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) expression() {}

type Function struct {
	Token  token.Token
	Params []*Identifier
	Body   []Node
}

func (f *Function) TokenLiteral() string {
	return f.Token.Literal
}

func (f *Function) expression() {}

type FunctionCall struct {
	Token token.Token
	FunctionExpr Expression
	Arguments []Expression
}

func (f *FunctionCall) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionCall) expression() {}

type String struct {
	Token token.Token
	Value string
}

func (s *String) TokenLiteral() string {
	return fmt.Sprintf("\"%s\"", s.Token.Literal)
}

func (s *String) expression() {}

type Array struct {
	Token token.Token
	Elements []Expression
}

func (a *Array) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Array) expression() {}
