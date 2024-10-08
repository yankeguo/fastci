package fastjs

import (
	"github.com/robertkrimen/otto"
)

// Object creates a new object with the given map values.
func Object(rp RuntimeProvider, m map[string]any) (obj *otto.Object, err error) {
	if obj, err = rp.Runtime().Object(`({})`); err != nil {
		return
	}
	for k, v := range m {
		if err = obj.Set(k, v); err != nil {
			return
		}
	}
	return
}

// Array creates a new array with the given slice values.
func Array[T any](rp RuntimeProvider, arr []T) (obj *otto.Object, err error) {
	if obj, err = rp.Runtime().Object(`([])`); err != nil {
		return
	}
	for _, v := range arr {
		if _, err = obj.Call("push", v); err != nil {
			return
		}
	}
	return
}
