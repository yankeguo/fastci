package fastci

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestCleanEnvKey(t *testing.T) {
	require.Equal(t, "HELLO_WORLD", cleanEnvKey("hello.world"))
	require.Equal(t, "HELLO__WORLD", cleanEnvKey("hello+.World"))
}

func TestToYaml(t *testing.T) {
	var data = struct {
		Hello string `json:"hello" yaml:"hello"`
	}{Hello: "world"}
	buf, err := json.Marshal(data)
	require.NoError(t, err)

	out, err := toYaml(buf)
	require.NoError(t, err)

	data.Hello = ""
	err = yaml.Unmarshal(out, &data)
	require.NoError(t, err)

	require.Equal(t, "world", data.Hello)

	buf = []byte("hello: world")
	out, err = toYaml(buf)
	require.NoError(t, err)
	require.Equal(t, buf, out)

	buf = []byte(":")
	out, err = toYaml(buf)
	require.Error(t, err)
	require.Empty(t, out)
}
