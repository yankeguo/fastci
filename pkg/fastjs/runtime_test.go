package fastjs

import (
	"os"
	"testing"

	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/require"
)

type runtime struct {
	*otto.Otto
}

func (r *runtime) Runtime() *otto.Otto {
	return r.Otto
}

func TestCreateEnvironObject(t *testing.T) {
	os.Setenv("HELLO", "")
	vm := otto.New()
	rp := &runtime{vm}
	obj, err := CreateEnvironObject(rp)
	require.NoError(t, err)
	val, err := obj.Get("PATH")
	require.NoError(t, err)
	require.True(t, val.IsDefined())
	require.True(t, val.IsString())
	val, err = obj.Get("HELLO")
	require.NoError(t, err)
	require.True(t, val.IsDefined())
	require.True(t, val.IsString())
	valStr, err := val.ToString()
	require.NoError(t, err)
	require.Equal(t, "", valStr)
}
