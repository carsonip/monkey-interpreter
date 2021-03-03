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

type Env struct {
	env map[string]object.Object
}

func NewEnv() *Env {
	env := &Env{}
	env.env = make(map[string]object.Object)
	return env
}

func NewEvaluator(parser *parser.Parser) Evaluator {
	return Evaluator{parser: parser, env: NewEnv()}
}

func (ev *Evaluator) EvalNext(env *Env) object.Object {
	node := ev.parser.NextNode()
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
	default:
		panic("not implemented")
	}
}

func (ev *Evaluator) evalLetStatement(statement *ast.LetStatement, env *Env) {
	name := statement.Name.TokenLiteral()
	val := ev.evalExpression(statement.Value, env)
	env.env[name] = val
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
		if obj, ok := env.env[expr.TokenLiteral()]; !ok {
			panic("unknown identifier")
		} else {
			return obj
		}
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
