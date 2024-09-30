package fastci

import (
	"context"
	"os"
	"strings"

	"github.com/robertkrimen/otto"
)

type Pipeline struct {
	Env map[string]string
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

func (p *Pipeline) setupEnv(vm *otto.Otto) {
	vm.Set("env", p.Env)
}

func (p *Pipeline) useDeployer1(call otto.FunctionCall) otto.Value {
	return otto.NullValue()
}

func (p *Pipeline) useDeployer2(call otto.FunctionCall) otto.Value {
	return otto.NullValue()
}

func (p *Pipeline) useJenkins(call otto.FunctionCall) otto.Value {
	return otto.NullValue()
}

func (p *Pipeline) setupFunctions(vm *otto.Otto) {
	vm.Set("useDeployer1", p.useDeployer1)
	vm.Set("deployer1", p.useDeployer1)
	vm.Set("useDeployer2", p.useDeployer2)
	vm.Set("deployer2", p.useDeployer2)
	vm.Set("useJenkins", p.useJenkins)
	vm.Set("jenkins", p.useJenkins)
}

func (p *Pipeline) Setup(vm *otto.Otto) {
	p.setupEnv(vm)
	p.setupFunctions(vm)
}

func (p *Pipeline) Do(ctx context.Context, script string) (err error) {
	vm := otto.New()
	p.Setup(vm)
	_, err = vm.Run(script)
	return
}
