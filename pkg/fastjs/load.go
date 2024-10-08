package fastjs

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

func LoadBoolField(out *bool, obj *otto.Object, name string) (err error) {
	var val otto.Value
	if val, err = obj.Get(name); err != nil {
		return
	}
	if val.IsUndefined() {
		return
	}
	if val.IsNull() {
		*out = false
		return
	}
	if *out, err = val.ToBoolean(); err != nil {
		return
	}
	return
}

func LoadStringField(out *string, obj *otto.Object, name string) (err error) {
	var val otto.Value
	if val, err = obj.Get(name); err != nil {
		return
	}
	if val.IsUndefined() {
		return
	}
	if val.IsNull() {
		*out = ""
		return
	}
	*out = val.String()
	return
}

func LoadFunctionField(out *otto.Value, obj *otto.Object, name string) (err error) {
	var val otto.Value
	if val, err = obj.Get(name); err != nil {
		return
	}
	if val.IsUndefined() {
		return
	}
	if val.IsNull() {
		*out = otto.Value{}
		return
	}
	if !val.IsObject() || val.Class() != "Function" {
		err = fmt.Errorf("field %s should be a function", name)
		return
	}
	*out = val
	return
}
