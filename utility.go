package fastci

import (
	"encoding/json"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	nonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

func SanitizeEnvName(name string) string {
	name = strings.ToUpper(name)
	name = nonAlphanumeric.ReplaceAllString(name, "_")
	return name
}

func ConvertJSONToYAML(buf []byte) (out []byte, err error) {
	var obj any
	if err = json.Unmarshal(buf, &obj); err != nil {
		return
	}
	out, err = yaml.Marshal(obj)
	return
}
