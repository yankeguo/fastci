package fastci

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/rg"
)

const (
	scriptBasic = `
console.log("\
HOME=\
"+ useEnv().HOME)
	`
)

func TestRunnerBasic(t *testing.T) {
	r := NewRunner()
	err := r.Execute(context.Background(), scriptBasic)
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

func TestRunnerTemporaryFiles(t *testing.T) {
	r := NewRunner()
	err := r.Execute(context.Background(), scriptTemp)
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

func TestRunnerScript(t *testing.T) {
	r := NewRunner()
	err := r.Execute(context.Background(), scriptScript)
	require.NoError(t, err)
}

const (
	scriptEnv = `
	useEnv()
	useEnv("hello", "World")
	console.log(useEnv("hello"))
	useEnv().hello = "world"
	`
)

func TestRunnerEnv(t *testing.T) {
	p := NewRunner()
	err := p.Execute(context.Background(), scriptEnv)
	require.NoError(t, err)
	require.Equal(t, "world", rg.Must(p.env.Get("hello")).String())
}
