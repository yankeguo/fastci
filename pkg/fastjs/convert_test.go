package fastjs

import (
	"testing"

	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/require"
)

func TestPlainObject(t *testing.T) {
	vm := otto.New()
	rp := &runtime{vm}
	obj, err := PlainObject(rp, map[string]interface{}{"a": 1})
	require.NoError(t, err)
	val, err := obj.Object().Get("a")
	require.NoError(t, err)
	valInt64, err := val.ToInteger()
	require.NoError(t, err)
	require.Equal(t, int64(1), valInt64)
}

func TestUnmarshalPlainObject(t *testing.T) {
	vm := otto.New()
	obj, err := vm.Object("({a:true})")
	require.NoError(t, err)
	var out map[string]interface{}
	err = UnmarshalPlainObject(&out, obj)
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{"a": true}, out)
}
