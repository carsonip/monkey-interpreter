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
	eval := NewEvaluator(&p, object.NewEnv())
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

func TestEvaluator_evalPrefixExpression(t *testing.T) {
	tests := [][]string{
		{"+1", "1"},
		{"-1", "-1"},
		{"!true", "false"},
		{"!false", "true"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalPrefixExpression_Error(t *testing.T) {
	tests := [][]string{
		{`+true`, "error: unsupported prefix operator on type"},
		{`+"foo"`, "error: unsupported prefix operator on type"},
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

func TestEvaluator_evalIdentifier_Error(t *testing.T) {
	tests := [][]string{
		{"x", "error: unknown identifier"},
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
		{"fn(){return fn(){return 1;}}()()", "1"},
	}
	runTests(t, tests)
}

func TestEvaluator_TestEvaluator_evalFunctionCall_Error(t *testing.T) {
	tests := [][]string{
		{"fn(){x; 1;}()", "error: unknown identifier"},
		{"fn(x){}()", "error: argument length mismatch"},
		{"fn(){}(1)", "error: argument length mismatch"},
		{"1(1)", "error: not a function"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalBuiltinFunction(t *testing.T) {
	tests := [][]string{
		{`len`, "builtin"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalFunctionCall_Scope(t *testing.T) {
	tests := [][]string{
		{"fn(){let x=1; fn(){let x = 2;}(); return x;}()", "1"},
		{"fn(){let x=1; return fn(x){return x;}(x+1);}()", "2"},
		{"let x=1; let f=fn(){let x=2; return fn(){return x;}}(); f();", "", "", "2"},
		{"let x=1; let f=fn(x){return fn(){return x;}}(2); f();", "", "", "2"},
		{"let x=1; let f=fn(){return x;}; x=2; f();", "", "", "2", "2"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalIfStatement(t *testing.T) {
	tests := [][]string{
		{"let x = 1; if (true) {x=2;}; x", "", "", "2"},
		{"let x = 1; if (false) {x=2;}; x", "", "", "1"},
		{"let x = 1; if (false) {x=2;} else {x=3;}; x", "", "", "3"},
		{"let x = 1; if (0) {x=2;} else {x=3;}; x", "", "", "2"},
		{"let x = 1; if (1) {x=2;} else {x=3;}; x", "", "", "2"},
		{"let x = 1; if (fn(){}) {x=2;} else {x=3;}; x", "", "", "2"},
		{"let x = 1; if (fn(){}()) {x=2;} else {x=3;}; x", "", "", "3"},
		{"let x = 1; if (true) {let x=2;}; x", "", "", "1"},
		{"fn(){if (true) {return 1; 2;}; return 3;}()", "1"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalIfStatement_Error(t *testing.T) {
	tests := [][]string{
		{"let x = 1; if (true) {y; x=2;}; x", "", "error: unknown identifier", "1"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalArithmetic(t *testing.T) {
	tests := [][]string{
		{"1 + 2", "3"},
		{"2 - 1", "1"},
		{"2 * 2", "4"},
		{"4 / 2", "2"},
		{"4 / 3", "1"},
		{`"a" + "b"`, `"ab"`},
	}
	runTests(t, tests)
}

func TestEvaluator_evalArithmetic_Error(t *testing.T) {
	tests := [][]string{
		{`"foo" * 2`, "error: unsupported types for arithmetic"},
		{`"foo" - "bar"`, "error: unsupported arithmetic operator"},
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

func TestEvaluator_evalComparison_Error(t *testing.T) {
	tests := [][]string{
		{"1 < true", "error: unsupported types for comparison"},
		{"[] > []", "error: unsupported types for comparison"},
		{"true > false", "error: unsupported comparison operator on type"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalAssignment(t *testing.T) {
	tests := [][]string{
		{"let x = 1; x = 2;", "", "2"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalAssignment_Error(t *testing.T) {
	tests := [][]string{
		{"1 = 1;", "error: bad lvalue"},
		{`"foo" = 1;`, "error: bad lvalue"},
	}
	runTests(t, tests)
}

func TestEvaluator_evalAssignmentIndex(t *testing.T) {
	tests := [][]string{
		{"let x = [1, 2]; x[1] = 3; x", "", "3", "[1, 3]"},
		{`let x = {"foo": "bar"}; x["foo"] = "baz"; x`, "", `"baz"`, `{"foo": "baz"}`},
		{`let x = {}; x["foo"] = "baz"; x`, "", `"baz"`, `{"foo": "baz"}`},
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

func TestEvaluator_Map(t *testing.T) {
	tests := [][]string{
		{`{1+1: 2+2}`, `{2: 4}`},
		{`{0: 1, false: 2}`, `{0: 1, false: 2}`},
		{`{"foo": {1: 2}}`, `{"foo": {1: 2}}`},
		{`{1: 1, 1: 2}`, `{1: 2}`},
	}
	runTests(t, tests)
}

func TestEvaluator_Index(t *testing.T) {
	tests := [][]string{
		{`[0][0]`, "0"},
		{`[0,1,1+1][1+1]`, "2"},
		{`[[1]][0][0]`, "1"},
		{`{1: 2}[1]`, "2"},
		{`{"foo": 2}["foo"]`, "2"},
		{`{0: 1, false: 2}[false]`, "2"},
		{`{0: 1, false: 2}[0]`, "1"},
		{`{0: {"foo": "bar"}}[0]["foo"]`, `"bar"`},
	}
	runTests(t, tests)
}

func TestEvaluator_Index_Error(t *testing.T) {
	tests := [][]string{
		{`1[0]`, "error: invalid type for index operation"},
		{`[1]["foo"];`, "error: array index not an integer"},
		{`[1][[]]`, "error: array index not an integer"},
		{`[][0]`, "error: array index out of bounds"},
		{`[][-1]`, "error: array index out of bounds"},
		{`[][0]=1`, "error: array index out of bounds"},
		{`{"foo": "bar"}["baz"]`, "error: key not found"},
		{`{}[[]]`, "error: key not hashable"},
		{`{}[[]]=1`, "error: key not hashable"},
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
