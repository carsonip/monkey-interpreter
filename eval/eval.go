package eval

import (
	"github.com/carsonip/monkey-interpreter/ast"
	"github.com/carsonip/monkey-interpreter/object"
	"github.com/carsonip/monkey-interpreter/parser"
	"github.com/carsonip/monkey-interpreter/token"
)

type Evaluator struct {
	parser *parser.Parser
}

func NewEvaluator(parser *parser.Parser) Evaluator {
	return Evaluator{parser: parser}
}

func (ev *Evaluator) EvalNext() object.Object {
	node := ev.parser.NextNode()
	if statement, ok := node.(ast.Statement); ok {
		ev.evalStatement(statement)
		return object.NULL
	} else if expr, ok := node.(ast.Expression); ok {
		return ev.evalExpression(expr)
	}
	panic("not implemented")
}

func (ev *Evaluator) evalStatement(statement ast.Statement) {

}

func (ev *Evaluator) evalExpression(expr ast.Expression) object.Object {
	switch expr := expr.(type) {
	case *ast.NumberLiteral:
		return object.NewInteger(expr.Value)
	case *ast.Boolean:
		return object.NewBoolean(expr.Value)
	case *ast.InfixExpression:
		return ev.evalInfixExpression(expr)
	case *ast.PrefixExpression:
		return ev.evalPrefixExpression(expr)
	}
	panic("not implemented")
}

func (ev *Evaluator) evalNumber(expr ast.Expression) int {
	obj := ev.evalExpression(expr)
	if num, ok := obj.(object.Integer); ok {
		return num.Value
	}
	panic("not int")
}

func (ev *Evaluator) evalBoolean(expr ast.Expression) bool {
	obj := ev.evalExpression(expr)
	if boolean, ok := obj.(object.Boolean); ok {
		return boolean.Value
	}
	panic("not bool")
}

func (ev *Evaluator) evalInfixExpression(infix *ast.InfixExpression) object.Object {
	var result int
	switch infix.Token.Type {
	case token.TOKEN_PLUS:
		result = ev.evalNumber(infix.Left) + ev.evalNumber(infix.Right)
	case token.TOKEN_MINUS:
		result = ev.evalNumber(infix.Left) - ev.evalNumber(infix.Right)
	case token.TOKEN_ASTERISK:
		result = ev.evalNumber(infix.Left) * ev.evalNumber(infix.Right)
	case token.TOKEN_SLASH:
		result = ev.evalNumber(infix.Left) / ev.evalNumber(infix.Right)
	}
	return object.NewInteger(result)
}

func (ev *Evaluator) evalPrefixExpression(prefix *ast.PrefixExpression) object.Object {
	switch prefix.Right.(type) {
	case *ast.NumberLiteral:
		var result int
		switch prefix.Token.Type {
		case token.TOKEN_PLUS:
			result = ev.evalNumber(prefix.Right)
		case token.TOKEN_MINUS:
			result = -ev.evalNumber(prefix.Right)
		default:
			panic("bad prefix")
		}
		return object.NewInteger(result)
	case *ast.Boolean:
		switch prefix.Token.Type {
		case token.TOKEN_NOT:
			return object.NewBoolean(!ev.evalBoolean(prefix.Right))
		}
	}
	panic("not implemented")
}
