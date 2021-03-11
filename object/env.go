package object

type Env struct {
	parentEnv   *Env
	env         map[string]Object
	returnValue Object
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
		env:       make(map[string]Object),
	}
	return env
}

func (e *Env) Get(name string) (Object, bool) {
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

func (e *Env) MustGet(name string) Object {
	val, ok := e.Get(name)
	if !ok {
		panic("unknown identifier")
	}
	return val
}

func (e *Env) SetNew(name string, value Object) {
	e.env[name] = value
}

func (e *Env) Set(name string, value Object) {
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

func (e *Env) Return(value Object) {
	e.returnValue = value
}

func (e *Env) Returned() (Object, bool) {
	if e.returnValue != nil {
		return e.returnValue, true
	}
	return nil, false
}

func (e *Env) MustReturned() Object {
	if returnValue, ok := e.Returned(); ok {
		return returnValue
	}
	return nil
}
