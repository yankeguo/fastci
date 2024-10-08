package fastci

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestSanitizeEnvKey(t *testing.T) {
	require.Equal(t, "HELLO_WORLD", SanitizeEnvName("hello.world"))
	require.Equal(t, "HELLO__WORLD", SanitizeEnvName("hello+.World"))
}

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

	buf = []byte("hello: world")
	out, err = ConvertJSONToYAML(buf)
	require.NoError(t, err)
	require.Equal(t, buf, out)

	buf = []byte(":")
	out, err = ConvertJSONToYAML(buf)
	require.Error(t, err)
	require.Empty(t, out)
}
