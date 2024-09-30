package fastci

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipeline(t *testing.T) {
	p := NewPipeline()
	err := p.Do(context.Background(), "console.log('HOME='+env.HOME)")
	require.NoError(t, err)
}
