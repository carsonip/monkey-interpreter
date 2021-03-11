package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_GetSet(t *testing.T) {
	env := NewEnv()
	val, ok := env.Get("foo")
	assert.False(t, ok)
	assert.Nil(t, val)
	obj := NewInteger(1)
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
	obj := NewInteger(1)
	env.Return(obj)
	val, ok = env.Returned()
	assert.True(t, ok)
	assert.Equal(t, obj, val)
}

func TestEnv_Nested(t *testing.T) {
	rootEnv := NewEnv()
	rootEnv.SetNew("foo", NewInteger(1))
	rootEnv.SetNew("bar", NewInteger(2))
	rootEnv.Return(NewInteger(3))
	env := NewNestedEnv(rootEnv)
	env.SetNew("bar", NewInteger(20))
	env.Return(NewInteger(30))

	assert.Equal(t, NewInteger(1), env.MustGet("foo"))
	assert.Equal(t, NewInteger(2), rootEnv.MustGet("bar"))
	assert.Equal(t, NewInteger(20), env.MustGet("bar"))
	assert.Equal(t, NewInteger(3), rootEnv.MustReturned())
	assert.Equal(t, NewInteger(30), env.MustReturned())
}
