package object

import (
	"fmt"
	"github.com/carsonip/monkey-interpreter/ast"
)

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
}

func (f Function) String() string {
	return "fn"
}

func NewFunction(params []string, body []ast.Node) Function {
	return Function{Params: params, Body: body}
}

type String struct {
	Value string
}

func (s String) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

func NewString(value string) String {
	return String{Value: value}
}
