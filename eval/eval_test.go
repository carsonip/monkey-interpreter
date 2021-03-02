package eval

import (
	"github.com/carsonip/monkey-interpreter/object"
	"github.com/carsonip/monkey-interpreter/parser"
	"github.com/carsonip/monkey-interpreter/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getEvaluator(input string) Evaluator {
	lexer := token.NewLexer(input)
	p := parser.NewParser(&lexer)
	eval := NewEvaluator(&p)
	return eval
}

func TestEvaluator_evalInfixExpression(t *testing.T) {
	eval := getEvaluator(`1+2*3-4`)
	obj := eval.EvalNext()
	num, ok := obj.(object.Integer)
	assert.True(t, ok)
	assert.Equal(t, 3, num.Value)
}

func TestEvaluator_evalPrefixExpression(t *testing.T) {
	tests := [][2]string{
		{"+1", "1"},
		{"-1", "-1"},
		{"!true", "false"},
		{"!false", "true"},
	}
	for _, test := range tests {
		eval := getEvaluator(test[0])
		obj := eval.EvalNext()
		assert.Equal(t, test[1], obj.String())
	}
}
