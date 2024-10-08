package fastci

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	nonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

func cleanEnvKey(name string) string {
	return nonAlphaNumeric.ReplaceAllString(strings.ToUpper(name), "_")
}

func toYaml(buf []byte) (out []byte, err error) {
	// test if it's already a valid yaml
	{
		var obj any
		if err = yaml.Unmarshal(buf, &obj); err == nil {
			out = buf
			return
		}
	}

	// test if it's a valid JSON
	{
		var obj any
		if err = json.Unmarshal(buf, &obj); err == nil {
			out, err = yaml.Marshal(obj)
			return
		}
	}

	err = errors.New("invalid input, cannot be converted to YAML")

	return
}
