package fastci

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/rg"
)

func runnerForTest(t *testing.T, script string) *Runner {
	r := NewRunner()
	r.noClear = true
	err := r.Execute(context.Background(), script)
	require.NoError(t, err)
	return r
}

func clearRunnerForTest(_ *testing.T, r *Runner) {
	r.noClear = false
	r.clear()
}

func TestRunnerBasic(t *testing.T) {
	runnerForTest(t, `
console.log("\
HOME=\
"+ useEnv().HOME)
	`)
}

func TestRunnerEnv(t *testing.T) {
	r := runnerForTest(t, `
	useEnv('hello', 'world')
	useEnv().hello = 'World'
	useEnv('hello', useEnv('hello')+"!")
	`)
	require.Equal(t, "World!", rg.Must(r.env.Get("hello")).String())
}

func TestRunnerRegistry(t *testing.T) {
	r := runnerForTest(t, `useRegistry('hello');useRegistry(useRegistry()+"!")`)
	require.Equal(t, "hello!", r.registry)
}

func TestRunnerImage(t *testing.T) {
	r := runnerForTest(t, `useImage('hello');useImage(useImage()+"!")`)
	require.Equal(t, "hello!", r.image)
}

func TestRunnerProfile(t *testing.T) {
	r := runnerForTest(t, `useProfile('hello');useProfile(useProfile()+"!")`)
	require.Equal(t, "hello!", r.profile)
}

func TestRunnerVersion(t *testing.T) {
	r := runnerForTest(t, `useVersion('hello');useVersion(useVersion()+"!")`)
	require.Equal(t, "hello!", r.version)
}

func TestRunnerDockerConfig(t *testing.T) {
	r := runnerForTest(t, `useDockerConfig(
	{
	content: {
	auths:{}}
	}
	)`)
	defer clearRunnerForTest(t, r)
	buf := rg.Must(os.ReadFile(filepath.Join(r.dockerConfig, "config.json")))
	require.Equal(t, `{"auths":{}}`, string(buf))
}

func TestRunnerKubeconfig(t *testing.T) {
	r := runnerForTest(t, `useKubeconfig(
	{
	content: {
	hello:'world'}
})`)
	defer clearRunnerForTest(t, r)
	buf := rg.Must(os.ReadFile(r.kubeconfig))
	require.Equal(t, "hello: world\n", string(buf))
}

func TestRunnerScript(t *testing.T) {
	runnerForTest(t, `useScript(
	'echo hello',
	'sleep 1',
	'echo world'
	);runScript();`)
}

func TestRunnerKubernetesWorkload(t *testing.T) {
	r := runnerForTest(t, `
	useKubernetesWorkload({namespace:'hello'})
	useKubernetesWorkload({name:'world'})
	useKubernetesWorkload({container:'nginx'})
	useKubernetesWorkload({kind:'Deployment'})
	useKubernetesWorkload({container:null})
	useKubernetesWorkload({init:1})
	`)
	require.Equal(t, "hello", r.workloadNamespace)
	require.Equal(t, "world", r.workloadName)
	require.Equal(t, "", r.workloadContainer)
	require.Equal(t, "Deployment", r.workloadKind)
	require.True(t, r.workloadInit)
}

func TestRunnerResolveCodingCredentials(t *testing.T) {
	r := runnerForTest(t, `
	useEnv('CODING_WEAVIN_USERNAME', 'hello')
	useEnv('CODING_WEAVIN_PASSWORD', 'foo')
	useEnv('CODING_WEAVIN_INFRA_PASSWORD', 'world')
	`)
	r.codingValuesTeam = "weavin"
	r.codingValuesProject = "infra"
	username, password := r.resolveCodingCredentials()
	require.Equal(t, "hello", username)
	require.Equal(t, "world", password)
}
