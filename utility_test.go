package fastci

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeEnvKey(t *testing.T) {
	require.Equal(t, "HELLO_WORLD", sanitizeEnvKey("hello.world"))
	require.Equal(t, "HELLO__WORLD", sanitizeEnvKey("hello+.World"))
}
