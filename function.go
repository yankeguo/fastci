package fastci

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/yankeguo/rg"
)

func CreateFunctionObjectGetSet(out *otto.Object, name string) func(call otto.FunctionCall) otto.Value {
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
				log.Printf("set env: %s", key.String())
			} else {
				panic("set env failed, key should be string")
			}
			return val
		}
	}
}

func CreateFunctionGetSetStringSlice(out *[]string, name string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		var newValues []string
		for _, val := range call.ArgumentList {
			newValues = append(newValues, val.String())
		}
		if len(newValues) > 0 {
			*out = newValues
			log.Printf("use %s: %v", name, *out)
		}
		return rg.Must(otto.ToValue(*out))
	}
}

func CreateFunctionGetSetString(out *string, name string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if val := call.Argument(0); val.IsString() {
			*out = val.String()
			log.Printf("use %s: %s", name, *out)
		}
		return rg.Must(otto.ToValue(*out))
	}
}

func CreateFunctionGetSetPathOrContent(out *string, name string, convert func(buf []byte) (out string, err error)) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		var (
			newContent []byte
			newPath    string
		)

		if val := call.Argument(0); val.IsString() {
			newContent = []byte(val.String())
		} else if val.IsObject() {
			buf := rg.Must(val.Object().MarshalJSON())

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
						if err := json.Unmarshal(data.Content, &s); err == nil {
							// string
							newContent = []byte(s)
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
			*out = rg.Must(convert(newContent))
			log.Println("use", name, "from content")
		}

		return rg.Must(otto.ToValue(*out))
	}
}
