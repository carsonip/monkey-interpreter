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
	eval := NewEvaluator(&p, NewEnv())
	return eval
}

func TestEvaluator_evalInfixExpression(t *testing.T) {
	eval := getEvaluator(`1+2*3-4`)
	obj := eval.EvalNext(eval.env)
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
		obj := eval.EvalNext(eval.env)
		assert.Equal(t, test[1], obj.String())
	}
}

func TestEvaluator_evalLetStatement(t *testing.T) {
	input := `let x = 100;
x`
	eval := getEvaluator(input)
	assert.Equal(t, "", eval.EvalNext(eval.env).String())
	assert.Equal(t, "100", eval.EvalNext(eval.env).String())
}

func TestEvaluator_evalIdentifier(t *testing.T) {
	input := `let x = 100;
x;
x+x;
3*x`
	eval := getEvaluator(input)
	assert.Equal(t, "", eval.EvalNext(eval.env).String())
	assert.Equal(t, "100", eval.EvalNext(eval.env).String())
	assert.Equal(t, "200", eval.EvalNext(eval.env).String())
	assert.Equal(t, "300", eval.EvalNext(eval.env).String())
}

func TestEvaluator_evalFunction(t *testing.T) {
	input := `fn(x, y){100; x+200;}`
	eval := getEvaluator(input)
	assert.Equal(t, "fn", eval.EvalNext(eval.env).String())
}

func TestEvaluator_evalFunctionCall(t *testing.T) {
	tests := [][2]string{
		{"fn(){1;}()", ""},
		{"fn(){1; return 2;}()", "2"},
		{"fn(x){1; return 2; return true;}(100)", "2"},
		{"fn(x, y){100; x+200; return x+y; 300;}(1, 2)", "3"},
		{"fn(){fn(){return 1;}()}()", ""},
		{"fn(){return fn(){return 1;}()}()", "1"},
	}
	for _, test := range tests {
		eval := getEvaluator(test[0])
		assert.Equal(t, test[1], eval.EvalNext(eval.env).String())
	}
}

func TestEvaluator_evalFunctionCall_Scope(t *testing.T) {
	tests := [][2]string{
		{"fn(){let x=1; fn(){let x = 2;}(); return x;}()", "1"},
		{"fn(){let x=1; return fn(x){return x;}(x+1);}()", "2"},
	}
	for _, test := range tests {
		eval := getEvaluator(test[0])
		assert.Equal(t, test[1], eval.EvalNext(eval.env).String())
	}
}

func TestEvaluator_evalIfStatement(t *testing.T) {
	tests := [][2]string{
		{"if true {let x=1;} else {let x=2;}; x", "1"},
		{"if false {let x=1;} else {let x=2;}; x", "2"},
		{"if true {}; fn(){if true {return 1; 2;}}()", "1"},
	}
	for _, test := range tests {
		eval := getEvaluator(test[0])
		assert.Equal(t, "", eval.EvalNext(eval.env).String())
		assert.Equal(t, test[1], eval.EvalNext(eval.env).String())
	}
}

func TestEvaluator_evalComparison(t *testing.T) {
	tests := [][2]string{
		{"1 < 2", "true"},
		{"2 > 1", "true"},
		{"1 + 0 < 2", "true"},
		{"1 + 1 < 2", "false"},
		{"1 + 1 > 2", "false"},
		{"2 > 1 + 0", "true"},
		{"1 + 1 == 2", "true"},
		{"1 != 2", "true"},
		{"true != false", "true"},
		{"true != true", "false"},
	}
	for _, test := range tests {
		eval := getEvaluator(test[0])
		assert.Equal(t, test[1], eval.EvalNext(eval.env).String())
	}
}
