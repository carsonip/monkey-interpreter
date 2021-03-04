package object

import (
	"fmt"
	"github.com/carsonip/monkey-interpreter/ast"
)

type Env struct {
	parentEnv *Env
	env map[string]Object
}

func NewEnv() *Env {
	env := &Env{
		env: make(map[string]Object),
	}
	return env
}

func NewNestedEnv(parentEnv *Env) *Env {
	env := &Env{
		parentEnv: parentEnv,
		env: make(map[string]Object),
	}
	return env
}

func (e *Env) Get(name string) Object {
	if obj, ok := e.env[name]; ok {
		return obj
	} else {
		if e.parentEnv != nil {
			return e.parentEnv.Get(name)
		} else {
			panic("unknown identifier")
		}
	}
}

func (e *Env) Set(name string, value Object) {
	e.env[name] = value
}

type Object interface {
	String() string
}

type Null struct {}

func (n Null) String() string {
	return ""
}

var NULL = Null{}

type Integer struct {
	Value int
}

func (i Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func NewInteger(value int) Integer {
	return Integer{Value: value}
}

type Boolean struct {
	Value bool
}

func (b Boolean) String() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}

func NewBoolean(value bool) Boolean {
	return Boolean{Value: value}
}

type Function struct {
	Params []string
	Body []ast.Node
	Env *Env
}

func (f Function) String() string {
	return "fn"
}

func NewFunction(params []*ast.Identifier, body []ast.Node, env *Env) Function {
	var paramStrs []string
	for _, p := range params {
		paramStrs = append(paramStrs, p.TokenLiteral())
	}
	return Function{Params: paramStrs, Body: body, Env: env}
}

type ReturnValue struct {
	Value Object
}

func (r ReturnValue) String() string {
	return r.Value.String()
}

func NewReturnValue(value Object) Object {
	return ReturnValue{Value: value}
}
