package fastjs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/yankeguo/rg"
)

type Function = func(call otto.FunctionCall) otto.Value

func GetterSetterForObject(rp RuntimeProvider, out *otto.Object, name string) Function {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 0 {
			return out.Value()
		} else if len(call.ArgumentList) == 1 {
			key := call.Argument(0).String()
			return rg.Must(out.Get(key))
		} else {
			key, val := call.Argument(0), call.Argument(1)
			if key.IsString() {
				out.Set(key.String(), val)
				log.Printf("set %s: %s", name, key.String())
			} else {
				rg.Must0(fmt.Errorf("set %s.%s failed, key should be string", name, key.String()))
			}
			return val
		}
	}
}

func GetterSetterForStringSlice(rp RuntimeProvider, out *[]string, name string) Function {
	return func(call otto.FunctionCall) otto.Value {
		var (
			values = []string{}
			update bool
		)

		if first := call.Argument(0); first.IsNull() {
			update = true
		} else if first.IsObject() && first.Class() == "Array" {
			buf := rg.Must(first.Object().MarshalJSON())
			rg.Must0(json.Unmarshal(buf, &values))
			update = true
		} else {
			for _, val := range call.ArgumentList {
				values = append(values, val.String())
			}
		}

		if update || len(values) > 0 {
			*out = values
			log.Printf("use %s: [%s]", name, strings.Join(*out, ", "))
		}

		return rg.Must(Array(rp, *out)).Value()
	}
}

func GetterSetterForString(rp RuntimeProvider, out *string, name string) Function {
	return func(call otto.FunctionCall) otto.Value {
		if first := call.Argument(0); first.IsString() {
			*out = first.String()
			log.Printf("use %s: %s", name, *out)
		}
		return rg.Must(otto.ToValue(*out))
	}
}

type LongStringPersister func(buf []byte, name string) (path string, err error)

func GetterSetterForLongString(rp RuntimeProvider, out *string, name string, persister LongStringPersister) Function {
	return func(call otto.FunctionCall) otto.Value {
		var (
			newContent []byte
			newPath    string
		)

		if first := call.Argument(0); first.IsString() {
			for i, val := range call.ArgumentList {
				if i > 0 {
					newContent = append(newContent, '\n')
				}
				newContent = append(newContent, []byte(val.String())...)
			}
		} else if first.IsObject() {
			buf := rg.Must(first.Object().MarshalJSON())

			var (
				lines []string
				data  struct {
					Content json.RawMessage `json:"content"`
					Base64  string          `json:"base64"`
					Path    string          `json:"path"`
				}
			)

			if err := json.Unmarshal(buf, &lines); err == nil {
				newContent = []byte(strings.Join(lines, "\n"))
			} else if err = json.Unmarshal(buf, &data); err == nil {
				if data.Path != "" {
					newPath = data.Path
				} else {
					if len(data.Content) > 0 {
						var s string
						var lines []string
						if err := json.Unmarshal(data.Content, &s); err == nil {
							// string
							newContent = []byte(s)
						} else if err := json.Unmarshal(data.Content, &lines); err == nil {
							// array of string
							newContent = []byte(strings.Join(lines, "\n"))
						} else {
							// object (raw)
							newContent = data.Content
						}
					} else if data.Base64 != "" {
						// base64
						newContent = rg.Must(base64.StdEncoding.DecodeString(data.Base64))
					}
				}
			}
		}

		if newPath != "" {
			*out = newPath
			log.Println("use", name, "from", newPath)
		} else if len(newContent) > 0 {
			*out = rg.Must(persister(newContent, name))
			log.Println("use", name, "from content")
		}

		return rg.Must(otto.ToValue(*out))
	}
}
