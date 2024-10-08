package fastjs

import (
	"testing"

	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/require"
)

func TestGetterSetterForObject(t *testing.T) {
	vm := otto.New()
	rp := &runtime{vm}
	obj, err := vm.Object("({hello:'world'})")
	require.NoError(t, err)
	vm.Set("fn", GetterSetterForObject(rp, obj, "obj"))

	val, err := vm.Eval("fn().hello")
	require.NoError(t, err)
	require.Equal(t, "world", val.String())

	_, err = vm.Eval("fn('hello', 'World')")
	require.NoError(t, err)

	val, err = vm.Eval("fn('hello')")
	require.NoError(t, err)
	require.Equal(t, "World", val.String())
}

func TestGetterSetterForStringSlice(t *testing.T) {
	vm := otto.New()
	rp := &runtime{vm}
	obj := []string{"hello"}
	vm.Set("fn", GetterSetterForStringSlice(rp, &obj, "obj"))

	val, err := vm.Eval("fn()[0]")
	require.NoError(t, err)
	require.Equal(t, "hello", val.String())

	_, err = vm.Eval("fn('hello', 'World')")
	require.NoError(t, err)

	val, err = vm.Eval("fn()[1]")
	require.NoError(t, err)
	require.Equal(t, "World", val.String())

	_, err = vm.Eval("fn([])")
	require.NoError(t, err)

	val, err = vm.Eval("String(fn().length)")
	require.NoError(t, err)
	require.Equal(t, "0", val.String())

	_, err = vm.Eval("fn(['Hello', 'world'])")
	require.NoError(t, err)

	val, err = vm.Eval("fn()[0]")
	require.NoError(t, err)
	require.Equal(t, "Hello", val.String())

	_, err = vm.Eval("fn(null)")
	require.NoError(t, err)

	val, err = vm.Eval("String(fn().length)")
	require.NoError(t, err)
	require.Equal(t, "0", val.String())
}

func TestGetterSetterForString(t *testing.T) {
	vm := otto.New()
	rp := &runtime{vm}
	obj := "hello"
	vm.Set("fn", GetterSetterForString(rp, &obj, "obj"))

	val, err := vm.Eval("fn()")
	require.NoError(t, err)
	require.Equal(t, "hello", val.String())

	_, err = vm.Eval("fn('World')")
	require.NoError(t, err)

	val, err = vm.Eval("fn()")
	require.NoError(t, err)
	require.Equal(t, "World", val.String())
}

func TestGetterSetterForLongString(t *testing.T) {
	vm := otto.New()
	rp := &runtime{vm}
	out := ""

	var (
		gotBuf  []byte
		gotName string
		retPath string
	)

	vm.Set("fn", GetterSetterForLongString(rp, &out, "obj", func(buf []byte, name string) (path string, err error) {
		gotBuf = buf
		gotName = name
		return retPath, nil
	}))

	out = "hello"

	val, err := vm.Eval("fn()")
	require.NoError(t, err)
	require.Equal(t, "hello", val.String())

	retPath = "path1"
	_, err = vm.Eval("fn('hello','world')")
	require.NoError(t, err)
	require.Equal(t, "path1", out)
	require.Equal(t, "hello\nworld", string(gotBuf))
	require.Equal(t, "obj", gotName)

	retPath = "path2"
	_, err = vm.Eval("fn({content:['hello','world']})")
	require.NoError(t, err)
	require.Equal(t, "path2", out)
	require.Equal(t, "hello\nworld", string(gotBuf))
	require.Equal(t, "obj", gotName)

	retPath = "path3"
	_, err = vm.Eval("fn({content:'hello world'})")
	require.NoError(t, err)
	require.Equal(t, "path3", out)
	require.Equal(t, "hello world", string(gotBuf))
	require.Equal(t, "obj", gotName)

	retPath = "path4"
	_, err = vm.Eval("fn({content:{hello:'world'}})")
	require.NoError(t, err)
	require.Equal(t, "path4", out)
	require.Equal(t, "{\"hello\":\"world\"}", string(gotBuf))
	require.Equal(t, "obj", gotName)

	retPath = "path5"
	gotBuf = []byte("not_set")
	gotName = "not_set"
	_, err = vm.Eval("fn({path:\"path6\"})")
	require.NoError(t, err)
	require.Equal(t, "path6", out)
	require.Equal(t, "not_set", string(gotBuf))
	require.Equal(t, "not_set", gotName)

	retPath = "path7"
	_, err = vm.Eval("fn({base64:'eyJoZWxsbyI6IndvcmxkIn0='})")
	require.NoError(t, err)
	require.Equal(t, "path7", out)
	require.Equal(t, "{\"hello\":\"world\"}", string(gotBuf))
	require.Equal(t, "obj", gotName)
}
