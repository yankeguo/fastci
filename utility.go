package fastci

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	regexpNonEnvName = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

// SanitizeEnvName sanitizes environment variable name by converting to uppercase and replacing non-alphanumeric characters with underscore
func SanitizeEnvName(name string) string {
	name = strings.ToUpper(name)
	name = regexpNonEnvName.ReplaceAllString(name, "_")
	return name
}

// ConvertJSONToYAML converts JSON to YAML if not already YAML
func ConvertJSONToYAML(buf []byte) (out []byte, err error) {
	var obj any
	if err = json.Unmarshal(buf, &obj); err != nil {
		// test if it's already YAML
		if err = yaml.Unmarshal(buf, &obj); err == nil {
			out = buf
		} else {
			err = errors.New("invalid JSON or YAML")
		}
		return
	}
	out, err = yaml.Marshal(obj)
	return
}
