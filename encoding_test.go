package fastci

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestConvertJSONToYAML(t *testing.T) {
	var data = struct {
		Hello string `json:"hello" yaml:"hello"`
	}{Hello: "world"}
	buf, err := json.Marshal(data)
	require.NoError(t, err)

	out, err := ConvertJSONToYAML(buf)
	require.NoError(t, err)

	data.Hello = ""
	err = yaml.Unmarshal(out, &data)
	require.NoError(t, err)

	require.Equal(t, "world", data.Hello)
}
