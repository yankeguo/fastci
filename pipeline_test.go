package fastci

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/rg"
)

func TestPipeline(t *testing.T) {
	for _, entry := range rg.Must(os.ReadDir(filepath.Join("testdata", "pipelines"))) {
		t.Run(entry.Name(), func(t *testing.T) {
			script := rg.Must(os.ReadFile(filepath.Join("testdata", "pipelines", entry.Name())))
			err := NewPipeline().Do(context.Background(), string(script))
			require.NoError(t, err)
		})
	}
}
