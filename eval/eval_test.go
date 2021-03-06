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

func runTests(t *testing.T, tests [][]string) {
	for _, inputOutput := range tests {
		input := inputOutput[0]
		outputs := inputOutput[1:]
		eval := getEvaluator(input)
		for _, output := range outputs {
			assert.Equal(t, output, eval.EvalNext(eval.env).String())
		}
		assert.Nil(t, eval.EvalNext(eval.env))
	}
}

func TestEvaluator_evalInfixExpression(t *testing.T) {
	eval := getEvaluator(`1+2*3-4`)
	obj := eval.EvalNext(eval.env)
	num, ok := obj.(object.Integer)
	assert.True(t, ok)
	assert.Equal(t, 3, num.Value)
}

func TestEvaluator_evalPrefixExpression(t *testing.T) {
	tests := [][]string{
		{"+1", "1"},
		{"-1", "-1"},
		{"!true", "false"},
		{"!false", "true"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalLetStatement(t *testing.T) {
	tests := [][]string{
		{`let x = 100; x`, "", "100"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalIdentifier(t *testing.T) {
	tests := [][]string{
		{`let x = 100; x; x+x; 3*x`, "", "100", "200", "300"},
		{`let len = 100; len`, "", "100"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalFunction(t *testing.T) {
	tests := [][]string{
		{`fn(x, y){100; x+200;}`, "fn"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalFunctionCall(t *testing.T) {
	tests := [][]string{
		{"fn(){1;}()", ""},
		{"fn(){1; return 2;}()", "2"},
		{"fn(x){1; return 2; return true;}(100)", "2"},
		{"fn(x, y){100; x+200; return x+y; 300;}(1, 2)", "3"},
		{"fn(){fn(){return 1;}()}()", ""},
		{"fn(){return fn(){return 1;}()}()", "1"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalFunctionCall_Scope(t *testing.T) {
	tests := [][]string{
		{"fn(){let x=1; fn(){let x = 2;}(); return x;}()", "1"},
		{"fn(){let x=1; return fn(x){return x;}(x+1);}()", "2"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalIfStatement(t *testing.T) {
	tests := [][]string{
		{"let x = 1; if (true) {x=2;}; x", "", "", "2"},
		{"let x = 1; if (true) {let x=2;}; x", "", "", "1"},
		{"fn(){if (true) {return 1; 2;}; return 3;}()", "1"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalComparison(t *testing.T) {
	tests := [][]string{
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
		{`"foo" == "bar"`, "false"},
		{`"foo" == "foo"`, "true"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalAssignment(t *testing.T) {
	tests := [][]string{
		{"let x = 1; x = 2;", "", "2"},
	}
	runTests(t, tests)
}

func TestEvaluator_String(t *testing.T) {
	tests := [][]string{
		{`"foo"`, `"foo"`},
	}
	runTests(t, tests)
}

func TestEvaluator_Array(t *testing.T) {
	tests := [][]string{
		{`[1+1]`, `[2]`},
		{`[["foo"]]`, `[["foo"]]`},
	}
	runTests(t, tests)
}

func TestEvaluator_Builtin(t *testing.T) {
	tests := [][]string{
		{`len([1,2])`, `2`},
		{`len("foo")`, `3`},
	}
	runTests(t, tests)
}
