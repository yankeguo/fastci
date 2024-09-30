package fastci

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/yankeguo/rg"
)

type Pipeline struct {
	Env map[string]string

	TemporaryDirectories []string

	Registry string
	Image    string
	Profile  string
	Version  string

	ScriptFile  string
	ScriptShell string

	DockerConfig string
	Kubeconfig   string
}

func NewPipeline() *Pipeline {
	p := &Pipeline{
		Env: make(map[string]string),
	}
	// set default env
	for _, entry := range os.Environ() {
		splits := strings.SplitN(entry, "=", 2)
		if len(splits) == 2 {
			p.Env[splits[0]] = splits[1]
		} else if len(splits) == 1 {
			p.Env[splits[0]] = ""
		}
	}
	return p
}

func (p *Pipeline) useDeployer1(call otto.FunctionCall) otto.Value {
	// TODO: implement deployer1
	return otto.NullValue()
}

func (p *Pipeline) useDeployer2(call otto.FunctionCall) otto.Value {
	// TODO: implement deployer2
	return otto.NullValue()
}

func (p *Pipeline) runScript(call otto.FunctionCall) otto.Value {
	shell := p.ScriptShell
	if shell == "" {
		shell = "/bin/bash"
	}
	cmd := exec.Command(shell, p.ScriptFile)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	rg.Must0(cmd.Run())
	return otto.NullValue()
}

func (p *Pipeline) setupEnv(vm *otto.Otto) {
	vm.Set("env", p.Env)
}

func (p *Pipeline) setupFunctions(vm *otto.Otto) {
	vm.Set("useDeployer1", p.useDeployer1)
	vm.Set("useDeployer2", p.useDeployer2)
	vm.Set("useRegistry", CreateFunctionGetSetString(&p.Registry, "registry"))
	vm.Set("useImage", CreateFunctionGetSetString(&p.Image, "image"))
	vm.Set("useProfile", CreateFunctionGetSetString(&p.Profile, "profile"))
	vm.Set("useVersion", CreateFunctionGetSetString(&p.Version, "version"))
	vm.Set("useDockerConfig", CreateFunctionGetSetPathOrContent(&p.DockerConfig, "docker config", func(buf []byte) (out string, err error) {
		_, out, err = p.createTemporaryFile("config.json", bytes.TrimSpace(buf))
		return
	}))
	vm.Set("useKubeconfig", CreateFunctionGetSetPathOrContent(&p.Kubeconfig, "kubeconfig", func(buf []byte) (out string, err error) {
		buf = bytes.TrimSpace(buf)
		if bytes.HasPrefix(buf, []byte("{")) {
			if buf, err = ConvertJSONToYAML(buf); err != nil {
				return
			}
		}
		out, _, err = p.createTemporaryFile("kubeconfig.yaml", buf)
		return
	}))
	vm.Set("useScript", CreateFunctionGetSetPathOrContent(&p.ScriptFile, "script", func(buf []byte) (out string, err error) {
		out, _, err = p.createTemporaryFile("script.sh", bytes.TrimSpace(buf))
		return
	}))
	vm.Set("useScriptShell", CreateFunctionGetSetString(&p.ScriptShell, "script shell"))
	vm.Set("runScript", p.runScript)
}

func (p *Pipeline) Setup(vm *otto.Otto) {
	p.setupEnv(vm)
	p.setupFunctions(vm)
}

func (p *Pipeline) Cleanup() {
	for _, dir := range p.TemporaryDirectories {
		log.Println("remove temporary directory:", dir)
		os.RemoveAll(dir)
	}
}

func (p *Pipeline) createTemporaryDirectory() (dir string, err error) {
	defer rg.Guard(&err)
	dir = rg.Must(os.MkdirTemp("", "fastci-*-tmp"))
	p.TemporaryDirectories = append(p.TemporaryDirectories, dir)
	return
}

func (p *Pipeline) createTemporaryFile(filename string, content []byte) (file string, dir string, err error) {
	defer rg.Guard(&err)
	dir = rg.Must(p.createTemporaryDirectory())
	file = filepath.Join(dir, filename)
	rg.Must0(os.WriteFile(file, content, 0644))
	return
}

func (p *Pipeline) Do(ctx context.Context, script string) (err error) {
	defer rg.Guard(&err)

	vm := otto.New()

	p.Setup(vm)
	defer p.Cleanup()

	_, err = vm.Run(script)
	return
}
