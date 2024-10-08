package fastjs

import (
	"testing"

	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/require"
)

func TestLoadBoolField(t *testing.T) {
	var out bool
	vm := otto.New()

	obj, err := vm.Object("({a:true})")
	require.NoError(t, err)
	LoadBoolField(&out, obj, "a")
	require.True(t, out)

	obj, err = vm.Object("({a:null})")
	require.NoError(t, err)
	LoadBoolField(&out, obj, "a")
	require.False(t, out)

	obj, err = vm.Object("({a:true})")
	require.NoError(t, err)
	LoadBoolField(&out, obj, "a")
	require.True(t, out)

	obj, err = vm.Object("({})")
	require.NoError(t, err)
	LoadBoolField(&out, obj, "a")
	require.True(t, out)
}

func TestLoadStringField(t *testing.T) {
	var out string
	vm := otto.New()

	obj, err := vm.Object("({a:'b'})")
	require.NoError(t, err)
	LoadStringField(&out, obj, "a")
	require.Equal(t, "b", out)

	obj, err = vm.Object("({a:null})")
	require.NoError(t, err)
	LoadStringField(&out, obj, "a")
	require.Equal(t, "", out)

	obj, err = vm.Object("({a:'b'})")
	require.NoError(t, err)
	LoadStringField(&out, obj, "a")
	require.Equal(t, "b", out)

	obj, err = vm.Object("({})")
	require.NoError(t, err)
	LoadStringField(&out, obj, "a")
	require.Equal(t, "b", out)
}

func TestLoadFunctionField(t *testing.T) {
	var out otto.Value
	vm := otto.New()

	obj, err := vm.Object("({a:function(){}})")
	require.NoError(t, err)
	LoadFunctionField(&out, obj, "a")
	require.True(t, out.IsFunction())

	obj, err = vm.Object("({a:null})")
	require.NoError(t, err)
	LoadFunctionField(&out, obj, "a")
	require.True(t, out.IsUndefined())

	obj, err = vm.Object("({a:function(){}})")
	require.NoError(t, err)
	LoadFunctionField(&out, obj, "a")
	require.True(t, out.IsFunction())

	obj, err = vm.Object("({})")
	require.NoError(t, err)
	LoadFunctionField(&out, obj, "a")
	require.True(t, out.IsFunction())
}
