package object

import "fmt"

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