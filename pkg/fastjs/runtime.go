package fastjs

import (
	"os"
	"strings"

	"github.com/robertkrimen/otto"
)

// RuntimeProvider is an interface for providing an otto runtime
type RuntimeProvider interface {
	// Runtime returns an otto runtime
	Runtime() *otto.Otto
}

// CreateEnvironObject creates an otto object from os.Environ()
func CreateEnvironObject(rp RuntimeProvider) (obj *otto.Object, err error) {
	if obj, err = rp.Runtime().Object(`({})`); err != nil {
		return
	}
	for _, entry := range os.Environ() {
		splits := strings.SplitN(entry, "=", 2)
		if len(splits) == 2 {
			if err = obj.Set(splits[0], splits[1]); err != nil {
				return
			}
		} else if len(splits) == 1 {
			if err = obj.Set(splits[0], ""); err != nil {
				return
			}
		}
	}
	return
}
