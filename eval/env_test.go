package eval

import (
	"github.com/carsonip/monkey-interpreter/object"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_GetSet(t *testing.T) {
	env := NewEnv()
	val, ok := env.Get("foo")
	assert.False(t, ok)
	assert.Nil(t, val)
	obj := object.NewInteger(1)
	env.SetNew("foo", obj)
	val, ok = env.Get("foo")
	assert.True(t, ok)
	assert.Equal(t, obj, val)
}

func TestEnv_Return(t *testing.T) {
	env := NewEnv()
	val, ok := env.Returned()
	assert.False(t, ok)
	assert.Nil(t, val)
	obj := object.NewInteger(1)
	env.Return(obj)
	val, ok = env.Returned()
	assert.True(t, ok)
	assert.Equal(t, obj, val)
}

func TestEnv_Nested(t *testing.T) {
	rootEnv := NewEnv()
	rootEnv.SetNew("foo", object.NewInteger(1))
	rootEnv.SetNew("bar", object.NewInteger(2))
	rootEnv.Return(object.NewInteger(3))
	env := NewNestedEnv(rootEnv)
	env.SetNew("bar", object.NewInteger(20))
	env.Return(object.NewInteger(30))

	assert.Equal(t, object.NewInteger(1), env.MustGet("foo"))
	assert.Equal(t, object.NewInteger(2), rootEnv.MustGet("bar"))
	assert.Equal(t, object.NewInteger(20), env.MustGet("bar"))
	assert.Equal(t, object.NewInteger(3), rootEnv.MustReturned())
	assert.Equal(t, object.NewInteger(30), env.MustReturned())
}
