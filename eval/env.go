package eval

import "github.com/carsonip/monkey-interpreter/object"

type Env struct {
	parentEnv   *Env
	env         map[string]object.Object
	returnValue object.Object
}

func NewEnv() *Env {
	env := &Env{
		env: make(map[string]object.Object),
	}
	return env
}

func NewNestedEnv(parentEnv *Env) *Env {
	env := &Env{
		parentEnv: parentEnv,
		env:       make(map[string]object.Object),
	}
	return env
}

func (e *Env) Get(name string) (object.Object, bool) {
	if obj, ok := e.env[name]; ok {
		return obj, true
	} else {
		if e.parentEnv != nil {
			return e.parentEnv.Get(name)
		} else {
			return nil, false
		}
	}
}

func (e *Env) MustGet(name string) object.Object {
	val, ok := e.Get(name)
	if !ok {
		panic("unknown identifier")
	}
	return val
}

func (e *Env) SetNew(name string, value object.Object) {
	e.env[name] = value
}

func (e *Env) Set(name string, value object.Object) {
	if _, ok := e.env[name]; ok {
		e.env[name] = value
		return
	} else {
		if e.parentEnv != nil {
			e.parentEnv.Set(name, value)
			return
		} else {
			panic("unknown identifier")
		}
	}
}

func (e *Env) Return(value object.Object) {
	e.returnValue = value
}

func (e *Env) Returned() (object.Object, bool) {
	if e.returnValue != nil {
		return e.returnValue, true
	}
	return nil, false
}

func (e *Env) MustReturned() object.Object {
	if returnValue, ok := e.Returned(); ok {
		return returnValue
	}
	return nil
}
