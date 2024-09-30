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
	_ "github.com/robertkrimen/otto/underscore"
	"github.com/yankeguo/rg"
)

type Runner struct {
	vm *otto.Otto

	env *otto.Object

	noClear bool

	tempDirs []string

	registry string
	image    string
	profile  string
	version  string

	buildScriptFile string
	buildShell      []string

	packageDockerfileFile string
	packageContext        string

	dockerConfig string
	kubeconfig   string
}

func NewRunner() *Runner {
	return &Runner{}
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

func (r *Runner) useDeployer1(call otto.FunctionCall) otto.Value {
	// TODO: implement deployer1
	return otto.NullValue()
}

func (r *Runner) useDeployer2(call otto.FunctionCall) otto.Value {
	// TODO: implement deployer2
	return otto.NullValue()
}

func (r *Runner) doBuild(call otto.FunctionCall) otto.Value {
	shell := r.buildShell
	if len(shell) == 0 {
		shell = []string{"/bin/bash"}
	}

	f := rg.Must(os.OpenFile(r.buildScriptFile, os.O_RDONLY, 0))
	defer f.Close()

	cmd := exec.Command(shell[0], shell[1:]...)
	cmd.Env = rg.Must(r.createEnviron())
	cmd.Stdin = f
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	rg.Must0(cmd.Run())

	return otto.NullValue()
}

func (r *Runner) doPackage(call otto.FunctionCall) otto.Value {
	return otto.NullValue()
}

func (r *Runner) doPublish(call otto.FunctionCall) otto.Value {
	return otto.NullValue()
}

func (r *Runner) setup() (err error) {
	r.vm = otto.New()

	// setup env
	if r.env, err = r.vm.Object("({})"); err != nil {
		return
	}
	for _, entry := range os.Environ() {
		splits := strings.SplitN(entry, "=", 2)
		if len(splits) == 2 {
			if err = r.env.Set(splits[0], splits[1]); err != nil {
				return
			}
		} else if len(splits) == 1 {
			if err = r.env.Set(splits[0], ""); err != nil {
				return
			}
		}
	}

	r.vm.Set("useEnv", FunctionForObjectGetterSetter(r.env, "env"))
	r.vm.Set("useDeployer1", r.useDeployer1)
	r.vm.Set("useDeployer2", r.useDeployer2)
	r.vm.Set("useRegistry", FunctionForStringGetterSetter(&r.registry, "registry"))
	r.vm.Set("useImage", FunctionForStringGetterSetter(&r.image, "image"))
	r.vm.Set("useProfile", FunctionForStringGetterSetter(&r.profile, "profile"))
	r.vm.Set("useVersion", FunctionForStringGetterSetter(&r.version, "version"))
	r.vm.Set("useDockerConfig", FunctionForLongStringGetterSetter(&r.dockerConfig, "docker config", func(buf []byte) (out string, err error) {
		_, out, err = r.createTempFile("config.json", bytes.TrimSpace(buf))
		return
	}))
	r.vm.Set("useKubeconfig", FunctionForLongStringGetterSetter(&r.kubeconfig, "kubeconfig", func(buf []byte) (out string, err error) {
		buf = bytes.TrimSpace(buf)
		if bytes.HasPrefix(buf, []byte("{")) {
			if buf, err = ConvertJSONToYAML(buf); err != nil {
				return
			}
		}
		out, _, err = r.createTempFile("kubeconfig.yaml", buf)
		return
	}))
	r.vm.Set("useBuildScript", FunctionForLongStringGetterSetter(&r.buildScriptFile, "build script", func(buf []byte) (out string, err error) {
		out, _, err = r.createTempFile("build.sh", bytes.TrimSpace(buf))
		return
	}))
	r.vm.Set("useBuildShell", FunctionForStringSliceGetterSetter(&r.buildShell, "build script shell"))
	r.vm.Set("doBuild", r.doBuild)
	r.vm.Set("usePackageDockerfile", FunctionForStringGetterSetter(&r.packageDockerfileFile, "package dockerfile"))
	r.vm.Set("usePackageContext", FunctionForStringGetterSetter(&r.packageContext, "package context"))
	r.vm.Set("doPackage", r.doPackage)
	r.vm.Set("doPublish", r.doPublish)
	return
}

func (r *Runner) clear() {
	if r.noClear {
		return
	}
	for _, dir := range r.tempDirs {
		log.Println("remove temporary directory:", dir)
		os.RemoveAll(dir)
	}
	r.tempDirs = nil
	r.env = nil
	r.vm = nil
}

func (r *Runner) createTempDir() (dir string, err error) {
	defer rg.Guard(&err)
	dir = rg.Must(os.MkdirTemp("", "fastci-*-tmp"))
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
	defer r.clear()

	_, err = r.vm.Run(script)
	return
}
