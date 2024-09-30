package fastci

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

func ConvertJSONToYAML(buf []byte) (out []byte, err error) {
	var obj any
	if err = json.Unmarshal(buf, &obj); err != nil {
		return
	}
	out, err = yaml.Marshal(obj)
	return
}
