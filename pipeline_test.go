package fastci

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	scriptBasic = `
console.log("\
HOME=\
"+ env.HOME)
	`
)

func TestPipelineBasic(t *testing.T) {
	p := NewPipeline()
	err := p.Do(context.Background(), scriptBasic)
	require.NoError(t, err)
}

const (
	scriptTemp = `
useDockerConfig({
    content: {
        auths: {
            "https://index.docker.io/v1/": {}
        }
    }
})

if (!useDockerConfig()) {
    throw new Error('dockerconfig() failed')
}

useKubeconfig("\
apiVersion: 'v1'\n\
")

if (!useKubeconfig()) {
    throw new Error('kubeconfig() failed')
}
	`
)

func TestPipelineTemporaryFiles(t *testing.T) {
	p := NewPipeline()
	err := p.Do(context.Background(), scriptTemp)
	require.NoError(t, err)
}

const (
	scriptScript = `
useBuildScript([
	'echo hello',
	'sleep 1',
	'echo world',
]);
doBuild();
`
)

func TestPipelineScript(t *testing.T) {
	p := NewPipeline()
	err := p.Do(context.Background(), scriptScript)
	require.NoError(t, err)
}
