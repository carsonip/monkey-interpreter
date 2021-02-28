package parser

import (
	"github.com/carsonip/monkey-interpreter/ast"
	"github.com/carsonip/monkey-interpreter/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	str := `let x = 123`
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
