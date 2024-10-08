package fastci

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"github.com/yankeguo/fastci/pkg/fastjs"
	"github.com/yankeguo/rg"
)

type Runner struct {
	vm  *otto.Otto
	env *otto.Object

	skipClean bool
	tempDirs  []string

	state struct {
		shell []string

		script struct {
			path string
		}

		docker struct {
			images         []string
			configPath     string
			dockerfilePath string
			buildContext   string
			buildArg       *otto.Object
		}

		kubernetes struct {
			kubeconfigPath string

			workload struct {
				namespace string
				name      string
				kind      string
				container string
				init      bool
			}
		}

		coding struct {
			values struct {
				team    string
				project string
				repo    string
				branch  string
				file    string
				update  otto.Value
			}
		}
	}
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Runtime() *otto.Otto {
	return r.vm
}

func (r *Runner) createEnviron() (items []string, err error) {
	for _, key := range r.env.Keys() {
		var val otto.Value
		if val, err = r.env.Get(key); err != nil {
			return
		}
		items = append(items, key+"="+val.String())
	}
	return
}

/*
func (r *Runner) createImages() (items []string, err error) {
	if r.state.registry == "" {
		err = errors.New("registry is not set")
		return
	}
	if r.state.image == "" {
		err = errors.New("image is not set")
		return
	}
	if r.state.profile == "" {
		if r.state.version == "" {
			items = append(items, path.Join(r.state.registry, r.state.image))
		} else {
			items = append(items, path.Join(r.state.registry, r.state.image))
			items = append(items, path.Join(r.state.registry, r.state.image+":"+r.state.version))
		}
	} else {
		if r.state.version == "" {
			items = append(items, path.Join(r.state.registry, r.state.image+":"+r.state.profile))
		} else {
			items = append(items, path.Join(r.state.registry, r.state.image+":"+r.state.profile))
			items = append(items, path.Join(r.state.registry, r.state.image+":"+r.state.profile+"-"+r.state.version))
		}
	}
	return
}
*/

func (r *Runner) useDeployer1(call otto.FunctionCall) otto.Value {
	//TODO: implement deployer1
	return otto.NullValue()
}

func (r *Runner) useDeployer2(call otto.FunctionCall) otto.Value {
	//TODO: implement deployer2
	return otto.NullValue()
}

func (r *Runner) runScript(call otto.FunctionCall) otto.Value {
	shell := r.state.shell
	if len(shell) == 0 {
		shell = []string{"bash"}
	}

	buf := rg.Must(os.ReadFile(r.state.script.path))

	log.Println("run script:", r.state.script.path, "\n", string(buf))

	cmd := exec.Command(shell[0], shell[1:]...)
	cmd.Env = rg.Must(r.createEnviron())
	cmd.Stdin = bytes.NewReader(buf)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	rg.Must0(cmd.Run())

	return otto.NullValue()
}

func (r *Runner) runDockerBuild(call otto.FunctionCall) otto.Value {
	if len(r.state.docker.images) == 0 {
		rg.Must0(errors.New("no images to build"))
		return otto.UndefinedValue()
	}
	var args []string

	// config
	if r.state.docker.configPath != "" {
		args = append(args, "--config", r.state.docker.configPath)
	}

	// command
	args = append(args, "buildx", "build", "--load")

	// build args
	for _, key := range r.state.docker.buildArg.Keys() {
		val := rg.Must(rg.Must(r.state.docker.buildArg.Get(key)).ToString())
		args = append(args, "--build-arg", key+"="+val)
	}

	// images
	for _, image := range r.state.docker.images {
		args = append(args, "-t", image)
	}

	// dockerfile
	if r.state.docker.dockerfilePath != "" {
		args = append(args, "-f", r.state.docker.dockerfilePath)
	}

	// build context
	if r.state.docker.buildContext != "" {
		args = append(args, r.state.docker.buildContext)
	} else {
		args = append(args, ".")
	}

	log.Println("run docker build:", strings.Join(args, " "))

	cmd := exec.Command("docker", args...)
	cmd.Env = rg.Must(r.createEnviron())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	rg.Must0(cmd.Run())

	return rg.Must(fastjs.PlainObject(r, r.state.docker.images))
}

func (r *Runner) runDockerPush(call otto.FunctionCall) otto.Value {
	if len(r.state.docker.images) == 0 {
		rg.Must0(errors.New("no images to push"))
		return otto.UndefinedValue()
	}
	for _, image := range r.state.docker.images {
		var args []string

		// config
		if r.state.docker.configPath != "" {
			args = append(args, "--config", r.state.docker.configPath)
		}

		// command
		args = append(args, "push", image)

		log.Println("run docker push:", strings.Join(args, " "))

		cmd := exec.Command("docker", args...)
		cmd.Env = rg.Must(r.createEnviron())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		rg.Must0(cmd.Run())
	}
	return rg.Must(fastjs.PlainObject(r, r.state.docker.images))
}

func (r *Runner) useKubernetesWorkload(call otto.FunctionCall) otto.Value {
	if arg := call.Argument(0); arg.IsObject() {
		obj := arg.Object()
		rg.Must0(fastjs.LoadStringField(&r.state.kubernetes.workload.namespace, obj, "namespace"))
		rg.Must0(fastjs.LoadStringField(&r.state.kubernetes.workload.name, obj, "name"))
		rg.Must0(fastjs.LoadStringField(&r.state.kubernetes.workload.kind, obj, "kind"))
		rg.Must0(fastjs.LoadStringField(&r.state.kubernetes.workload.container, obj, "container"))
		rg.Must0(fastjs.LoadBoolField(&r.state.kubernetes.workload.init, obj, "init"))
	}
	return rg.Must(fastjs.PlainObject(r, map[string]any{
		"namespace": r.state.kubernetes.workload.namespace,
		"name":      r.state.kubernetes.workload.name,
		"kind":      r.state.kubernetes.workload.kind,
		"container": r.state.kubernetes.workload.container,
		"init":      r.state.kubernetes.workload.init,
	}))
}

func (r *Runner) resolveCodingCredentials() (username string, password string) {
	var parts []string
	if r.state.coding.values.team != "" {
		parts = append(parts, SanitizeEnvName(r.state.coding.values.team))
		if r.state.coding.values.project != "" {
			parts = append(parts, SanitizeEnvName(r.state.coding.values.project))
			if r.state.coding.values.repo != "" {
				parts = append(parts, SanitizeEnvName(r.state.coding.values.repo))
			}
		}
	}
	usernameKeys := []string{"CODING_USERNAME"}
	passwordKeys := []string{"CODING_PASSWORD"}
	for i := range parts {
		usernameKeys = append(usernameKeys, "CODING_"+strings.Join(parts[:i+1], "_")+"_USERNAME")
		passwordKeys = append(passwordKeys, "CODING_"+strings.Join(parts[:i+1], "_")+"_PASSWORD")
	}
	slices.Reverse(usernameKeys)
	slices.Reverse(passwordKeys)

	for _, usernameKey := range usernameKeys {
		val := rg.Must(r.env.Get(usernameKey))
		if val.IsString() {
			username = val.String()
			log.Println("use coding username from:", usernameKey)
			break
		}
	}

	for _, passwordKey := range passwordKeys {
		val := rg.Must(r.env.Get(passwordKey))
		if val.IsString() {
			password = val.String()
			log.Println("use coding password from:", passwordKey)
			break
		}
	}
	return
}

func (r *Runner) useCodingValues(call otto.FunctionCall) otto.Value {
	if arg := call.Argument(0); arg.IsObject() {
		obj := arg.Object()
		rg.Must0(fastjs.LoadStringField(&r.state.coding.values.team, obj, "team"))
		rg.Must0(fastjs.LoadStringField(&r.state.coding.values.project, obj, "project"))
		rg.Must0(fastjs.LoadStringField(&r.state.coding.values.repo, obj, "repo"))
		rg.Must0(fastjs.LoadStringField(&r.state.coding.values.branch, obj, "branch"))
		rg.Must0(fastjs.LoadStringField(&r.state.coding.values.file, obj, "file"))
		rg.Must0(fastjs.LoadFunctionField(&r.state.coding.values.update, obj, "update"))
	}
	return rg.Must(fastjs.PlainObject(r, map[string]any{
		"team":    r.state.coding.values.team,
		"project": r.state.coding.values.project,
		"repo":    r.state.coding.values.repo,
		"branch":  r.state.coding.values.branch,
		"file":    r.state.coding.values.file,
		"update":  r.state.coding.values.update,
	}))
}

func (r *Runner) deployKubernetesWorkload(call otto.FunctionCall) otto.Value {
	return otto.NullValue()
}

func (r *Runner) deployCodingValues(call otto.FunctionCall) otto.Value {
	return otto.NullValue()
}

func (r *Runner) setup() (err error) {
	r.vm = otto.New()

	if r.state.docker.buildArg, err = r.vm.Object("({})"); err != nil {
		return
	}

	if r.env, err = fastjs.CreateEnvironObject(r); err != nil {
		return
	}

	r.vm.Set("useShell", fastjs.GetterSetterForStringSlice(r, &r.state.shell, "shell"))
	r.vm.Set("useEnv", fastjs.GetterSetterForObject(r, r.env, "env"))
	r.vm.Set("useDockerConfig", fastjs.GetterSetterForLongString(r, &r.state.docker.configPath, "docker config", func(buf []byte, name string) (out string, err error) {
		_, out, err = r.createTempFile("config.json", bytes.TrimSpace(buf))
		return
	}))
	r.vm.Set("useKubeconfig", fastjs.GetterSetterForLongString(r, &r.state.kubernetes.kubeconfigPath, "kubeconfig", func(buf []byte, name string) (out string, err error) {
		buf = rg.Must(ConvertJSONToYAML(bytes.TrimSpace(buf)))
		out, _, err = r.createTempFile("kubeconfig.yaml", buf)
		return
	}))

	r.vm.Set("useScript", fastjs.GetterSetterForLongString(r, &r.state.script.path, "script", func(buf []byte, name string) (out string, err error) {
		out, _, err = r.createTempFile("script.sh", bytes.TrimSpace(buf))
		return
	}))
	r.vm.Set("runScript", r.runScript)

	r.vm.Set("useDockerImages", fastjs.GetterSetterForStringSlice(r, &r.state.docker.images, "docker images"))
	r.vm.Set("useDockerBuildArg", fastjs.GetterSetterForObject(r, r.state.docker.buildArg, "docker build arg"))
	r.vm.Set("useDockerfile", fastjs.GetterSetterForString(r, &r.state.docker.dockerfilePath, "dockerfile"))
	r.vm.Set("useDockerBuildContext", fastjs.GetterSetterForString(r, &r.state.docker.buildContext, "docker context"))
	r.vm.Set("runDockerBuild", r.runDockerBuild)
	r.vm.Set("runDockerPush", r.runDockerPush)

	r.vm.Set("useKubernetesWorkload", r.useKubernetesWorkload)
	r.vm.Set("deployKubernetesWorkload", r.deployKubernetesWorkload)

	r.vm.Set("useCodingValues", r.useCodingValues)
	r.vm.Set("deployCodingValues", r.deployCodingValues)

	r.vm.Set("useDeployer1", r.useDeployer1)
	r.vm.Set("useDeployer2", r.useDeployer2)
	return
}

func (r *Runner) clean() {
	if r.skipClean {
		return
	}
	for _, dir := range r.tempDirs {
		log.Println("remove temporary directory:", dir)
		os.RemoveAll(dir)
	}
	r.tempDirs = nil
	r.env = nil
	r.state.coding.values.update = otto.Value{} // reset to undefined
	r.vm = nil
}

func (r *Runner) createTempDir() (dir string, err error) {
	defer rg.Guard(&err)

	buf := make([]byte, 12)
	rg.Must(rand.Read(buf))

	dir = filepath.Join(os.TempDir(), "fastci-"+strconv.FormatInt(time.Now().Unix(), 10)+"-"+hex.EncodeToString(buf))
	rg.Must0(os.MkdirAll(dir, 0755))

	r.tempDirs = append(r.tempDirs, dir)
	return
}

func (r *Runner) createTempFile(filename string, content []byte) (file string, dir string, err error) {
	defer rg.Guard(&err)
	dir = rg.Must(r.createTempDir())
	file = filepath.Join(dir, filename)
	rg.Must0(os.WriteFile(file, content, 0644))
	return
}

func (r *Runner) Execute(ctx context.Context, script any) (err error) {
	defer rg.Guard(&err)

	rg.Must0(r.setup())
	defer r.clean()

	_, err = r.vm.Run(script)
	return
}
