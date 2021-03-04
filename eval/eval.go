package eval

import (
	"github.com/carsonip/monkey-interpreter/ast"
	"github.com/carsonip/monkey-interpreter/object"
	"github.com/carsonip/monkey-interpreter/parser"
	"github.com/carsonip/monkey-interpreter/token"
)

type Evaluator struct {
	parser *parser.Parser
	env *Env
}

func NewEvaluator(parser *parser.Parser) Evaluator {
	return Evaluator{parser: parser, env: NewEnv()}
}

func (ev *Evaluator) EvalNext(env *Env) object.Object {
	node := ev.parser.NextNode()
	return ev.Eval(node, env)
}

func (ev *Evaluator) Eval(node ast.Node, env *Env) object.Object {
	if statement, ok := node.(ast.Statement); ok {
		ev.evalStatement(statement, env)
		return object.NULL
	} else if expr, ok := node.(ast.Expression); ok {
		return ev.evalExpression(expr, env)
	}
	panic("not implemented")
}

func (ev *Evaluator) evalStatement(statement ast.Statement, env *Env) {
	switch statement := statement.(type) {
	case *ast.LetStatement:
		ev.evalLetStatement(statement, env)
	case *ast.ReturnStatement:
		ev.evalReturnStatement(statement, env)
	case *ast.IfStatement:
		ev.evalIfStatement(statement, env)
	default:
		panic("not implemented")
	}
}

func (ev *Evaluator) evalLetStatement(statement *ast.LetStatement, env *Env) {
	name := statement.Name.TokenLiteral()
	val := ev.evalExpression(statement.Value, env)
	env.Set(name, val)
}

func (ev *Evaluator) evalReturnStatement(statement *ast.ReturnStatement, env *Env) {
	val := ev.evalExpression(statement.Value, env)
	env.Return(val)
}

func (ev *Evaluator) evalIfStatement(statement *ast.IfStatement, env *Env) {
	var nodes []ast.Node
	pass := ev.evalBoolean(statement.Condition, env)
	if pass {
		nodes = statement.Then
	} else {
		nodes = statement.Else
	}
	for _, node := range nodes {
		ev.Eval(node, env)
		if _, ok := env.Returned(); ok {
			return
		}
	}
}

func (ev *Evaluator) evalExpression(expr ast.Expression, env *Env) object.Object {
	switch expr := expr.(type) {
	case *ast.NumberLiteral:
		return object.NewInteger(expr.Value)
	case *ast.Boolean:
		return object.NewBoolean(expr.Value)
	case *ast.InfixExpression:
		return ev.evalInfixExpression(expr, env)
	case *ast.PrefixExpression:
		return ev.evalPrefixExpression(expr, env)
	case *ast.Identifier:
		return env.MustGet(expr.TokenLiteral())
	case *ast.Function:
		return ev.evalFunction(expr, env)
	case *ast.FunctionCall:
		return ev.evalFunctionCall(expr, env)
	}
	panic("not implemented")
}

func (ev *Evaluator) evalNumber(expr ast.Expression, env *Env) int {
	obj := ev.evalExpression(expr, env)
	if num, ok := obj.(object.Integer); ok {
		return num.Value
	}
	panic("not int")
}

func (ev *Evaluator) evalBoolean(expr ast.Expression, env *Env) bool {
	obj := ev.evalExpression(expr, env)
	if boolean, ok := obj.(object.Boolean); ok {
		return boolean.Value
	}
	panic("not bool")
}

func (ev *Evaluator) evalInfixExpression(infix *ast.InfixExpression, env *Env) object.Object {
	var result int
	switch infix.Token.Type {
	case token.TOKEN_PLUS:
		result = ev.evalNumber(infix.Left, env) + ev.evalNumber(infix.Right, env)
	case token.TOKEN_MINUS:
		result = ev.evalNumber(infix.Left, env) - ev.evalNumber(infix.Right, env)
	case token.TOKEN_ASTERISK:
		result = ev.evalNumber(infix.Left, env) * ev.evalNumber(infix.Right, env)
	case token.TOKEN_SLASH:
		result = ev.evalNumber(infix.Left, env) / ev.evalNumber(infix.Right, env)
	}
	return object.NewInteger(result)
}

func (ev *Evaluator) evalPrefixExpression(prefix *ast.PrefixExpression, env *Env) object.Object {
	switch prefix.Right.(type) {
	case *ast.NumberLiteral:
		var result int
		switch prefix.Token.Type {
		case token.TOKEN_PLUS:
			result = ev.evalNumber(prefix.Right, env)
		case token.TOKEN_MINUS:
			result = -ev.evalNumber(prefix.Right, env)
		default:
			panic("bad prefix")
		}
		return object.NewInteger(result)
	case *ast.Boolean:
		switch prefix.Token.Type {
		case token.TOKEN_NOT:
			return object.NewBoolean(!ev.evalBoolean(prefix.Right, env))
		}
	}
	panic("not implemented")
}

func (ev *Evaluator) evalFunction(fn *ast.Function, env *Env) *object.Function {
	var params []string
	for _, p := range fn.Params {
		params = append(params, p.TokenLiteral())
	}
	fnObj := object.NewFunction(params, fn.Body)
	return &fnObj
}

func (ev *Evaluator) evalFunctionCall(fnCall *ast.FunctionCall, env *Env) object.Object {
	expr := ev.evalExpression(fnCall.FunctionExpr, env)
	if fn, ok := expr.(*object.Function); !ok {
		panic("not a function")
	} else {
		var args []object.Object
		for _, argExpr := range fnCall.Arguments {
			args = append(args, ev.evalExpression(argExpr, env))
		}
		return ev.callFunction(fn, args, env)
	}
}

func (ev *Evaluator) callFunction(fn *object.Function, args []object.Object, parentEnv *Env) object.Object {
	env := NewNestedEnv(parentEnv)
	if len(fn.Params) != len(args) {
		panic("argument length mismatch")
	}
	for i, name := range fn.Params {
		env.Set(name, args[i])
	}
	for _, node := range fn.Body {
		ev.Eval(node, env)
		if val, ok := env.Returned(); ok {
			return val
		}
	}
	return object.NULL
}
