package fastjs

import (
	"testing"

	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {
	vm := otto.New()
	rp := &runtime{vm}
	obj, err := Object(rp, map[string]interface{}{"a": 1})
	require.NoError(t, err)
	val, err := obj.Get("a")
	require.NoError(t, err)
	valInt64, err := val.ToInteger()
	require.NoError(t, err)
	require.Equal(t, int64(1), valInt64)
}

func TestArray(t *testing.T) {
	vm := otto.New()
	rp := &runtime{vm}
	obj, err := Array(rp, []int{1})
	require.NoError(t, err)
	val, err := obj.Get("length")
	require.NoError(t, err)
	valInt64, err := val.ToInteger()
	require.NoError(t, err)
	require.Equal(t, int64(1), valInt64)
}
