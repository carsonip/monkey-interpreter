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
	if num, ok := expr.(*ast.NumberLiteral); ok {
		return object.NewInteger(num.Value)
	} else if infix, ok := expr.(*ast.InfixExpression); ok {
		return ev.evalInfixExpression(infix)
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
