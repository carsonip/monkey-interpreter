package parser

import (
	"github.com/carsonip/monkey-interpreter/ast"
	"github.com/carsonip/monkey-interpreter/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_LetStatement(t *testing.T) {
	str := `let x = 123;`
	lex := token.NewLexer(str)
	p := NewParser(&lex)
	node := p.NextNode()
	l, ok := node.(*ast.LetStatement)
	assert.True(t, ok)
	assert.Equal(t, "x", l.Name.TokenLiteral())
	num, ok := l.Value.(*ast.NumberLiteral)
	assert.True(t, ok)
	assert.Equal(t, 123, num.Value)
}

func TestParser_InfixExpression(t *testing.T) {
	str := `1 + 2`
	lex := token.NewLexer(str)
	p := NewParser(&lex)
	node := p.NextNode()
	exp, ok := node.(*ast.InfixExpression)
	assert.True(t, ok)
	assert.Equal(t, token.TOKEN_PLUS, exp.Token.Type)
	num, ok := exp.Left.(*ast.NumberLiteral)
	assert.True(t, ok)
	assert.Equal(t, 1, num.Value)
	num, ok = exp.Right.(*ast.NumberLiteral)
	assert.True(t, ok)
	assert.Equal(t, 2, num.Value)
}

func TestParser_InfixExpression_Precedence(t *testing.T) {
	str := `1 + 2 * 3 - 4 / 5`
	lex := token.NewLexer(str)
	p := NewParser(&lex)
	node := p.NextNode()
	exp, ok := node.(*ast.InfixExpression)
	assert.True(t, ok)
	assert.Equal(t, token.TOKEN_MINUS, exp.Token.Type)
	lExp, ok := exp.Left.(*ast.InfixExpression)
	assert.True(t, ok)
	assert.Equal(t, token.TOKEN_PLUS, lExp.Token.Type)
	lRExp, ok := lExp.Right.(*ast.InfixExpression)
	assert.True(t, ok)
	assert.Equal(t, token.TOKEN_ASTERISK, lRExp.Token.Type)
	rExp, ok := exp.Right.(*ast.InfixExpression)
	assert.True(t, ok)
	assert.Equal(t, token.TOKEN_SLASH, rExp.Token.Type)
}

func TestParser_InfixExpression_Precedence_Left(t *testing.T) {
	str := `1 + 2 + 3`
	lex := token.NewLexer(str)
	p := NewParser(&lex)
	node := p.NextNode()
	exp, ok := node.(*ast.InfixExpression)
	assert.True(t, ok)
	assert.Equal(t, token.TOKEN_PLUS, exp.Token.Type)
	lExp, ok := exp.Left.(*ast.InfixExpression)
	assert.True(t, ok)
	assert.Equal(t, token.TOKEN_PLUS, lExp.Token.Type)
	lLExp, ok := lExp.Left.(*ast.NumberLiteral)
	assert.True(t, ok)
	assert.Equal(t, 1, lLExp.Value)
	lRExp, ok := lExp.Right.(*ast.NumberLiteral)
	assert.True(t, ok)
	assert.Equal(t, 2, lRExp.Value)
	rExp, ok := exp.Right.(*ast.NumberLiteral)
	assert.True(t, ok)
	assert.Equal(t, 3, rExp.Value)
}

func TestParser_PrefixExpression(t *testing.T) {
	str := `-1`
	lex := token.NewLexer(str)
	p := NewParser(&lex)
	node := p.NextNode()
	exp, ok := node.(*ast.PrefixExpression)
	assert.True(t, ok)
	assert.Equal(t, token.TOKEN_MINUS, exp.Token.Type)
	num, ok := exp.Right.(*ast.NumberLiteral)
	assert.True(t, ok)
	assert.Equal(t, 1, num.Value)
}
