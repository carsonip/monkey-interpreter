package eval

import (
	"github.com/carsonip/monkey-interpreter/object"
)

var BUILTINS = map[string]object.BuiltinFunction{
	"len": {Fn: _len},
}

func _len(args ...object.Object) object.Object {
	if len(args) != 1 {
		panic("bad args len for Len")
	}
	obj := args[0]
	switch obj := obj.(type) {
	case object.Array:
		return object.NewInteger(len(obj.Elements))
	case object.String:
		return object.NewInteger(len(obj.Value))
	default:
		panic("unsupported type for Len")
	}
}
