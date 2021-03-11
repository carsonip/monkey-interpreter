package eval

import (
	"github.com/carsonip/monkey-interpreter/ast"
	"github.com/carsonip/monkey-interpreter/object"
	"github.com/carsonip/monkey-interpreter/parser"
	"github.com/carsonip/monkey-interpreter/token"
)

type Evaluator struct {
	parser *parser.Parser
	env *object.Env
}

func NewEvaluator(parser *parser.Parser, env *object.Env) Evaluator {
	return Evaluator{parser: parser, env: env}
}

func (ev *Evaluator) EvalNext(env *object.Env) object.Object {
	node := ev.parser.NextNode()
	if node == nil {
		return nil
	}
	return ev.Eval(node, env)
}

func (ev *Evaluator) Eval(node ast.Node, env *object.Env) (ret object.Object) {
	defer func() {
		err := recover()
		if err != nil {
			if e, ok := err.(object.Error); ok {
				ret = e
			}
		}
	}()
	if statement, ok := node.(ast.Statement); ok {
		ev.evalStatement(statement, env)
		return object.NULL
	} else if expr, ok := node.(ast.Expression); ok {
		return ev.evalExpression(expr, env)
	}
	panic(object.NewError("not implemented"))
}

func (ev *Evaluator) evalStatement(statement ast.Statement, env *object.Env) {
	switch statement := statement.(type) {
	case *ast.LetStatement:
		ev.evalLetStatement(statement, env)
	case *ast.ReturnStatement:
		ev.evalReturnStatement(statement, env)
	case *ast.IfStatement:
		ev.evalIfStatement(statement, env)
	default:
		panic(object.NewError("not implemented"))
	}
}

func (ev *Evaluator) evalLetStatement(statement *ast.LetStatement, env *object.Env) {
	name := statement.Name.TokenLiteral()
	val := ev.evalExpression(statement.Value, env)
	env.SetNew(name, val)
}

func (ev *Evaluator) evalReturnStatement(statement *ast.ReturnStatement, env *object.Env) {
	val := ev.evalExpression(statement.Value, env)
	env.Return(val)
}

func (ev *Evaluator) evalIfStatement(statement *ast.IfStatement, env *object.Env) {
	var nodes []ast.Node
	ok := isTruthy(ev.evalExpression(statement.Condition, env))
	if ok {
		nodes = statement.Then
	} else {
		nodes = statement.Else
	}
	newEnv := object.NewNestedEnv(env)
	for _, node := range nodes {
		result := ev.Eval(node, newEnv)
		if returnValue, ok := newEnv.Returned(); ok {
			env.Return(returnValue)
			return
		} else if err, ok := result.(object.Error); ok {
			panic(err)
		}
	}
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case object.Boolean:
		return obj.Value
	case object.Null:
		return false
	default:
		return true
	}
}

func (ev *Evaluator) evalExpression(expr ast.Expression, env *object.Env) object.Object {
	switch expr := expr.(type) {
	case *ast.NumberLiteral:
		return object.NewInteger(expr.Value)
	case *ast.Boolean:
		return object.NewBoolean(expr.Value)
	case *ast.String:
		return object.NewString(expr.Value)
	case *ast.InfixExpression:
		return ev.evalInfixExpression(expr, env)
	case *ast.PrefixExpression:
		return ev.evalPrefixExpression(expr, env)
	case *ast.Identifier:
		return ev.evalIdentifier(expr, env)
	case *ast.Function:
		return ev.evalFunction(expr, env)
	case *ast.FunctionCall:
		return ev.evalFunctionCall(expr, env)
	case *ast.Array:
		return ev.evalArray(expr, env)
	case *ast.Index:
		return ev.evalIndex(expr, env)
	case *ast.Map:
		return ev.evalMap(expr, env)
	}
	panic(object.NewError("not implemented"))
}

func (ev *Evaluator) evalIdentifier(expr ast.Expression, env *object.Env) object.Object {
	name := expr.TokenLiteral()
	if val, ok := env.Get(name); ok {
		return val
	} else if builtin, ok := BUILTINS[name]; ok {
		return builtin
	} else {
		panic(object.NewError("unknown identifier"))
	}
}

func (ev *Evaluator) evalNumber(expr ast.Expression, env *object.Env) int {
	obj := ev.evalExpression(expr, env)
	if num, ok := obj.(object.Integer); ok {
		return num.Value
	}
	panic(object.NewError("not int"))
}

func (ev *Evaluator) evalBoolean(expr ast.Expression, env *object.Env) bool {
	obj := ev.evalExpression(expr, env)
	if boolean, ok := obj.(object.Boolean); ok {
		return boolean.Value
	}
	panic(object.NewError("not bool"))
}

func (ev *Evaluator) evalInfixExpression(infix *ast.InfixExpression, env *object.Env) object.Object {
	switch infix.Token.Type {
	case token.TOKEN_PLUS:
		return object.NewInteger(ev.evalNumber(infix.Left, env) + ev.evalNumber(infix.Right, env))
	case token.TOKEN_MINUS:
		return object.NewInteger(ev.evalNumber(infix.Left, env) - ev.evalNumber(infix.Right, env))
	case token.TOKEN_ASTERISK:
		return object.NewInteger(ev.evalNumber(infix.Left, env) * ev.evalNumber(infix.Right, env))
	case token.TOKEN_SLASH:
		return object.NewInteger(ev.evalNumber(infix.Left, env) / ev.evalNumber(infix.Right, env))
	case token.TOKEN_EQUAL, token.TOKEN_NOTEQUAL, token.TOKEN_LT, token.TOKEN_GT:
		return ev.evalComparison(infix.Left, infix.Right, infix.Token.Type, env)
	case token.TOKEN_ASSIGNMENT:
		return ev.evalAssignment(infix.Left, infix.Right, env)
	}
	panic(object.NewError("unknown infix operator type"))
}

func (ev *Evaluator) evalComparison(leftExpr ast.Expression, rightExpr ast.Expression, tokenType token.TokenType, env *object.Env) object.Object {
	left := ev.evalExpression(leftExpr, env)
	right := ev.evalExpression(rightExpr, env)
	switch left := left.(type) {
	case object.Integer:
		if right, ok := right.(object.Integer); ok {
			switch tokenType {
			case token.TOKEN_EQUAL:
				return object.NewBoolean(left.Value == right.Value)
			case token.TOKEN_NOTEQUAL:
				return object.NewBoolean(left.Value != right.Value)
			case token.TOKEN_LT:
				return object.NewBoolean(left.Value < right.Value)
			case token.TOKEN_GT:
				return object.NewBoolean(left.Value > right.Value)
			}
		} else {
			panic(object.NewError("comparison type mismatch"))
		}
	case object.Boolean:
		if right, ok := right.(object.Boolean); ok {
			switch tokenType {
			case token.TOKEN_EQUAL:
				return object.NewBoolean(left.Value == right.Value)
			case token.TOKEN_NOTEQUAL:
				return object.NewBoolean(left.Value != right.Value)
			default:
				panic(object.NewError("cannot compare boolean"))
			}
		} else {
			panic(object.NewError("comparison type mismatch"))
		}
	case object.String:
		if right, ok := right.(object.String); ok {
			switch tokenType {
			case token.TOKEN_EQUAL:
				return object.NewBoolean(left.Value == right.Value)
			case token.TOKEN_NOTEQUAL:
				return object.NewBoolean(left.Value != right.Value)
			default:
				panic(object.NewError("cannot compare string"))
			}
		} else {
			panic(object.NewError("comparison type mismatch"))
		}
	}
	panic(object.NewError("unknown type for comparison"))
}

func (ev *Evaluator) evalAssignment(left ast.Expression, right ast.Expression, env *object.Env) object.Object {
	val := ev.evalExpression(right, env)
	switch left := left.(type) {
	case *ast.Identifier:
		name := left.TokenLiteral()
		env.Set(name, val)
	case *ast.Index:
		ev.evalAssignmentIndex(left, right, env)
	default:
		panic(object.NewError("bad lvalue"))
	}
	return val
}

func (ev *Evaluator) evalAssignmentIndex(ind *ast.Index, right ast.Expression, env *object.Env) {
	left := ev.evalExpression(ind.Left, env)
	indVal := ev.evalExpression(ind.Index, env)
	value := ev.evalExpression(right, env)
	switch left := left.(type) {
	case object.Array:
		left.Set(indVal, value)
	case object.Map:
		left.Set(indVal, value)
	}
}

func (ev *Evaluator) evalPrefixExpression(prefix *ast.PrefixExpression, env *object.Env) object.Object {
	switch prefix.Right.(type) {
	case *ast.NumberLiteral:
		var result int
		switch prefix.Token.Type {
		case token.TOKEN_PLUS:
			result = ev.evalNumber(prefix.Right, env)
		case token.TOKEN_MINUS:
			result = -ev.evalNumber(prefix.Right, env)
		default:
			panic(object.NewError("bad prefix"))
		}
		return object.NewInteger(result)
	case *ast.Boolean:
		switch prefix.Token.Type {
		case token.TOKEN_NOT:
			return object.NewBoolean(!ev.evalBoolean(prefix.Right, env))
		}
	}
	panic(object.NewError("not implemented"))
}

func (ev *Evaluator) evalFunction(fn *ast.Function, env *object.Env) object.Function {
	var params []string
	for _, p := range fn.Params {
		params = append(params, p.TokenLiteral())
	}
	fnObj := object.NewFunction(params, fn.Body)
	return fnObj
}

func (ev *Evaluator) convertFnArgs(argExprs []ast.Expression, env *object.Env) []object.Object {
	var args []object.Object
	for _, argExpr := range argExprs {
		args = append(args, ev.evalExpression(argExpr, env))
	}
	return args
}

func (ev *Evaluator) evalFunctionCall(fnCall *ast.FunctionCall, env *object.Env) object.Object {
	expr := ev.evalExpression(fnCall.FunctionExpr, env)
	switch fn := expr.(type) {
	case object.Function:
		args := ev.convertFnArgs(fnCall.Arguments, env)
		return ev.callFunction(fn, args, env)
	case object.BuiltinFunction:
		args := ev.convertFnArgs(fnCall.Arguments, env)
		return ev.callBuiltinFunction(fn, args, env)
	default:
		panic(object.NewError("not a function"))
	}
}

func (ev *Evaluator) callFunction(fn object.Function, args []object.Object, parentEnv *object.Env) object.Object {
	env := object.NewNestedEnv(parentEnv)
	if len(fn.Params) != len(args) {
		panic(object.NewError("argument length mismatch"))
	}
	for i, name := range fn.Params {
		env.SetNew(name, args[i])
	}
	for _, node := range fn.Body {
		result := ev.Eval(node, env)
		if val, ok := env.Returned(); ok {
			return val
		} else if err, ok := result.(object.Error); ok {
			panic(err)
		}
	}
	return object.NULL
}

func (ev *Evaluator) callBuiltinFunction(fn object.BuiltinFunction, args []object.Object, parentEnv *object.Env) object.Object {
	return fn.Fn(args...)
}

func (ev *Evaluator) evalArray(arr *ast.Array, env *object.Env) object.Array {
	var elements []object.Object
	for _, expr := range arr.Elements {
		elements = append(elements, ev.evalExpression(expr, env))
	}
	arrObj := object.NewArray(elements)
	return arrObj
}

func (ev *Evaluator) evalIndex(ind *ast.Index, env *object.Env) object.Object {
	var obj object.Object
	left := ev.evalExpression(ind.Left, env)
	switch left := left.(type) {
	case object.Array:
		indNum := ev.evalNumber(ind.Index, env)
		if indNum < 0 || indNum >= len(left.Elements) {
			panic(object.NewError("bad index value"))
		}
		obj = left.Elements[indNum]
	case object.Map:
		key := ev.evalExpression(ind.Index, env)
		if val, ok := left.Get(key); !ok {
			panic(object.NewError("key not found"))
		} else {
			obj = val
		}
	default:
		panic(object.NewError("invalid type for index operation"))
	}
	return obj
}

func (ev *Evaluator) evalMap(m *ast.Map, env *object.Env) object.Map {
	var pairs [][2]object.Object
	for _, kvExprs := range m.Pairs {
		kExpr := kvExprs[0]
		vExpr := kvExprs[1]
		k := ev.evalExpression(kExpr, env)
		v := ev.evalExpression(vExpr, env)
		pairs = append(pairs, [2]object.Object{k, v})
	}
	return object.NewMap(pairs)
}
