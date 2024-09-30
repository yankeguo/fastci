package fastci

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/yankeguo/rg"
	"gopkg.in/yaml.v3"
)

type Pipeline struct {
	Env map[string]string

	TemporaryDirectories []string
	TemporaryFiles       []string

	Registry string
	Image    string
	Profile  string
	Version  string

	BuildDir         string
	BuildScriptShell string
	BuildScriptFile  string

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

func (p *Pipeline) useRegistry(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return rg.Must(otto.ToValue(p.Registry))
	}
	p.Registry = call.Argument(0).String()
	log.Println("use registry:", p.Registry)
	return otto.NullValue()
}

func (p *Pipeline) useJenkins(call otto.FunctionCall) otto.Value {
	// TODO: implement jenkins
	return otto.NullValue()
}

func (p *Pipeline) useDockerConfig(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return rg.Must(otto.ToValue(p.DockerConfig))
	}
	val := call.Argument(0)
	if val.IsObject() {
		buf := rg.Must(val.Object().MarshalJSON())
		dir := rg.Must(p.createTemporaryDirectory())
		rg.Must0(os.WriteFile(filepath.Join(dir, "config.json"), buf, 0640))
		p.DockerConfig = dir
		log.Println("use docker config:", p.DockerConfig)
	} else if val.IsString() {
		p.DockerConfig = val.String()
		log.Print("use docker config:", p.DockerConfig)
	}
	return otto.NullValue()
}

func (p *Pipeline) useKubeconfig(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return rg.Must(otto.ToValue(p.Kubeconfig))
	}
	val := call.Argument(0)
	if val.IsObject() {
		buf := rg.Must(val.Object().MarshalJSON())
		var m map[string]any
		rg.Must0(json.Unmarshal(buf, &m))
		buf = rg.Must(yaml.Marshal(m))
		p.Kubeconfig = rg.Must(p.createTemporaryFile(buf))
	} else if val.IsString() {
		p.Kubeconfig = val.String()
	}
	log.Println("use kubeconfig:", p.Kubeconfig)
	return otto.NullValue()
}

func (p *Pipeline) useImage(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return rg.Must(otto.ToValue(p.Image))
	}
	p.Image = call.Argument(0).String()
	log.Println("use image:", p.Image)
	return otto.NullValue()
}

func (p *Pipeline) useProfile(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return rg.Must(otto.ToValue(p.Profile))
	}
	p.Profile = call.Argument(0).String()
	log.Println("use profile:", p.Profile)
	return otto.NullValue()
}

func (p *Pipeline) useVersion(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return rg.Must(otto.ToValue(p.Version))
	}
	p.Version = call.Argument(0).String()
	log.Println("use version:", p.Version)
	return otto.NullValue()
}

func (p *Pipeline) useBuildScript(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return rg.Must(otto.ToValue(p.BuildScriptFile))
	}
	script := call.Argument(0).String()
	p.BuildScriptFile = rg.Must(p.createTemporaryFile([]byte(script)))
	log.Println("use build script:", script)
	return otto.NullValue()
}

func (p *Pipeline) useBuildScriptFile(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return rg.Must(otto.ToValue(p.BuildScriptFile))
	}
	file := call.Argument(0).String()
	p.BuildScriptFile = file
	log.Println("use build script file:", p.BuildScriptFile)
	return otto.NullValue()
}

func (p *Pipeline) useBuildScriptShell(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return rg.Must(otto.ToValue(p.BuildScriptShell))
	}
	p.BuildScriptShell = call.Argument(0).String()
	log.Println("use build script shell:", p.BuildScriptShell)
	return otto.NullValue()
}

func (p *Pipeline) useBuildDir(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		return rg.Must(otto.ToValue(p.BuildDir))
	}
	p.BuildDir = call.Argument(0).String()
	log.Println("use build dir:", p.BuildDir)
	return otto.NullValue()
}

func (p *Pipeline) doBuild(call otto.FunctionCall) otto.Value {
	//TODO: implement doBuild
	return otto.NullValue()
}

func (p *Pipeline) setupEnv(vm *otto.Otto) {
	vm.Set("env", p.Env)
}

func (p *Pipeline) setupFunctions(vm *otto.Otto) {
	vm.Set("useDeployer1", p.useDeployer1)
	vm.Set("useDeployer2", p.useDeployer2)
	vm.Set("useRegistry", p.useRegistry)
	vm.Set("useImage", p.useImage)
	vm.Set("useProfile", p.useProfile)
	vm.Set("useVersion", p.useVersion)
	vm.Set("useJenkins", p.useJenkins)
	vm.Set("useDockerConfig", p.useDockerConfig)
	vm.Set("useKubeconfig", p.useKubeconfig)
	vm.Set("useBuildScript", p.useBuildScript)
	vm.Set("useBuildScriptFile", p.useBuildScriptFile)
	vm.Set("useBuildScriptShell", p.useBuildScriptShell)
	vm.Set("useBuildDir", p.useBuildDir)
	vm.Set("doBuild", p.doBuild)
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
	for _, file := range p.TemporaryFiles {
		log.Println("remove temporary file:", file)
		os.Remove(file)
	}
}

func (p *Pipeline) createTemporaryDirectory() (dir string, err error) {
	defer rg.Guard(&err)
	dir = rg.Must(os.MkdirTemp("", "fastci-*-dockerconfig"))
	p.TemporaryDirectories = append(p.TemporaryDirectories, dir)
	return
}

func (p *Pipeline) createTemporaryFile(buf []byte) (file string, err error) {
	defer rg.Guard(&err)
	_file := rg.Must(os.CreateTemp("", "fastci-*-kubeconfig"))
	defer _file.Close()
	file = _file.Name()
	p.TemporaryFiles = append(p.TemporaryFiles, file)
	rg.Must(_file.Write(buf))
	return
}

func (p *Pipeline) Do(ctx context.Context, script string) (err error) {
	vm := otto.New()

	p.Setup(vm)
	defer p.Cleanup()

	_, err = vm.Run(script)
	return
}
