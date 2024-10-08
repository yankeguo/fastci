package fastjs

import (
	"encoding/json"

	"github.com/robertkrimen/otto"
)

func PlainObject(rp RuntimeProvider, val any) (out otto.Value, err error) {
	var buf []byte
	if buf, err = json.Marshal(val); err != nil {
		return
	}
	var obj *otto.Object
	if obj, err = rp.Runtime().Object("(" + string(buf) + ")"); err != nil {
		return
	}
	out = obj.Value()
	return
}

func UnmarshalPlainObject(out any, obj *otto.Object) (err error) {
	var buf []byte
	if buf, err = obj.MarshalJSON(); err != nil {
		return
	}
	if err = json.Unmarshal(buf, out); err != nil {
		return
	}
	return
}
