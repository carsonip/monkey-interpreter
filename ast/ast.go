package ast

import (
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
	return i.Token.Literal
}

func (i *InfixExpression) expression() {}
