package fastci

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
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

func (p *Pipeline) useStringVar(out *string, name string, call otto.FunctionCall) otto.Value {
	if val := call.Argument(0); val.IsString() {
		*out = val.String()
		log.Printf("use %s: %s", name, *out)
	}
	return rg.Must(otto.ToValue(*out))
}

func (p *Pipeline) useContentOrPathVar(out *string, name string, fnSave func(buf []byte) (out string, err error), call otto.FunctionCall) otto.Value {
	var (
		newContent []byte
		newPath    string
	)

	if val := call.Argument(0); val.IsString() {
		newContent = []byte(val.String())
	} else if val.IsObject() {
		buf := rg.Must(val.Object().MarshalJSON())

		var (
			lines []string
			data  struct {
				Content json.RawMessage `json:"content"`
				Base64  string          `json:"base64"`
				Path    string          `json:"path"`
			}
		)

		if err := json.Unmarshal(buf, &lines); err == nil {
			newContent = []byte(strings.Join(lines, "\n"))
		} else if err = json.Unmarshal(buf, &data); err == nil {
			if data.Path != "" {
				newPath = data.Path
			} else {
				if len(data.Content) > 0 {
					var s string
					if err := json.Unmarshal(data.Content, &s); err == nil {
						// string
						newContent = []byte(s)
					} else {
						// object (raw)
						newContent = data.Content
					}
				} else if data.Base64 != "" {
					// base64
					newContent = rg.Must(base64.StdEncoding.DecodeString(data.Base64))
				}
			}
		}
	}

	if newPath != "" {
		*out = newPath
		log.Println("use", name, "from", newPath)
	} else if len(newContent) > 0 {
		*out = rg.Must(fnSave(newContent))
		log.Println("use", name, "from content")
	}

	return rg.Must(otto.ToValue(*out))
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
	return p.useStringVar(&p.Registry, "registry", call)
}

func (p *Pipeline) useImage(call otto.FunctionCall) otto.Value {
	return p.useStringVar(&p.Image, "image", call)
}

func (p *Pipeline) useProfile(call otto.FunctionCall) otto.Value {
	return p.useStringVar(&p.Profile, "profile", call)
}

func (p *Pipeline) useVersion(call otto.FunctionCall) otto.Value {
	return p.useStringVar(&p.Version, "version", call)
}

func (p *Pipeline) useDockerConfig(call otto.FunctionCall) otto.Value {
	return p.useContentOrPathVar(&p.DockerConfig, "docker config", func(buf []byte) (out string, err error) {
		buf = bytes.TrimSpace(buf)
		if _, out, err = p.createTemporaryFile("config.json", buf); err != nil {
			return
		}
		return
	}, call)
}

func (p *Pipeline) useKubeconfig(call otto.FunctionCall) otto.Value {
	return p.useContentOrPathVar(&p.Kubeconfig, "kubeconfig", func(buf []byte) (out string, err error) {
		buf = bytes.TrimSpace(buf)
		if bytes.HasPrefix(buf, []byte("{")) {
			if buf, err = ConvertJSONToYAML(buf); err != nil {
				return
			}
		}
		if out, _, err = p.createTemporaryFile("kubeconfig.yaml", buf); err != nil {
			return
		}
		return
	}, call)
}

func (p *Pipeline) useScript(call otto.FunctionCall) otto.Value {
	return p.useContentOrPathVar(&p.ScriptFile, "script", func(buf []byte) (out string, err error) {
		buf = bytes.TrimSpace(buf)
		if out, _, err = p.createTemporaryFile("script.sh", buf); err != nil {
			return
		}
		return
	}, call)
}

func (p *Pipeline) useScriptShell(call otto.FunctionCall) otto.Value {
	return p.useStringVar(&p.ScriptShell, "script shell", call)
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
	vm.Set("useRegistry", p.useRegistry)
	vm.Set("useImage", p.useImage)
	vm.Set("useProfile", p.useProfile)
	vm.Set("useVersion", p.useVersion)
	vm.Set("useDockerConfig", p.useDockerConfig)
	vm.Set("useKubeconfig", p.useKubeconfig)
	vm.Set("useScript", p.useScript)
	vm.Set("useScriptShell", p.useScriptShell)
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
