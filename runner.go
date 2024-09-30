package fastci

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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

	script      string
	scriptShell []string

	dockerfile    string
	dockerContext string

	workloadNamespace string
	workloadName      string
	workloadKind      string
	workloadContainer string
	workloadInit      bool

	codingValuesTeam    string
	codingValuesProject string
	codingValuesRepo    string
	codingValuesBranch  string
	codingValuesFile    string
	codingValuesUpdate  otto.Value

	dockerConfig string
	kubeconfig   string
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) createGetterSetterForObject(out *otto.Object, name string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 0 {
			return out.Value()
		} else if len(call.ArgumentList) == 1 {
			key := call.Argument(0).String()
			return rg.Must(out.Get(key))
		} else {
			key, val := call.Argument(0), call.Argument(1)
			if key.IsString() {
				out.Set(key.String(), val)
				log.Printf("set %s: %s", name, key.String())
			} else {
				panic(fmt.Sprintf("set %s.%s failed, key should be string", name, key.String()))
			}
			return val
		}
	}
}

func (r *Runner) createGetterSetterForStringSlice(out *[]string, name string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		var newValues []string
		for _, val := range call.ArgumentList {
			newValues = append(newValues, val.String())
		}
		if len(newValues) > 0 {
			*out = newValues
			log.Printf("use %s: [%s]", name, strings.Join(*out, ", "))
		}
		return rg.Must(otto.ToValue(*out))
	}
}

func (r *Runner) createGetterSetterForString(out *string, name string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if val := call.Argument(0); val.IsString() {
			*out = val.String()
			log.Printf("use %s: %s", name, *out)
		}
		return rg.Must(otto.ToValue(*out))
	}
}

func (r *Runner) createGetterSetterForLongString(out *string, name string, convert func(buf []byte) (out string, err error)) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		var (
			newContent []byte
			newPath    string
		)

		if val := call.Argument(0); val.IsString() {
			for i, val := range call.ArgumentList {
				if i > 0 {
					newContent = append(newContent, '\n')
				}
				newContent = append(newContent, []byte(val.String())...)
			}
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
						var lines []string
						if err := json.Unmarshal(data.Content, &s); err == nil {
							// string
							newContent = []byte(s)
						} else if err := json.Unmarshal(data.Content, &lines); err == nil {
							// array of string
							newContent = []byte(strings.Join(lines, "\n"))
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
			*out = rg.Must(convert(newContent))
			log.Println("use", name, "from content")
		}

		return rg.Must(otto.ToValue(*out))
	}
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

func (r *Runner) createPlainObject(m map[string]any) (ret otto.Value, err error) {
	var obj *otto.Object
	if obj, err = r.vm.Object("({})"); err != nil {
		return
	}
	for key, val := range m {
		if err = obj.Set(key, val); err != nil {
			return
		}
	}
	ret = obj.Value()
	return
}

func (r *Runner) loadObjectBoolField(out *bool, obj *otto.Object, name string) {
	val := rg.Must(obj.Get(name))
	if val.IsUndefined() {
		return
	}
	if val.IsNull() {
		*out = false
		return
	}
	*out = rg.Must(val.ToBoolean())
}

func (r *Runner) loadObjectStringField(out *string, obj *otto.Object, name string) {
	val := rg.Must(obj.Get(name))
	if val.IsUndefined() {
		return
	}
	if val.IsNull() {
		*out = ""
		return
	}
	*out = val.String()
}

func (r *Runner) loadObjectFunctionField(out *otto.Value, obj *otto.Object, name string) {
	val := rg.Must(obj.Get(name))
	if val.IsUndefined() {
		return
	}
	if val.IsNull() {
		*out = otto.Value{}
		return
	}
	if val.IsObject() {
		if val.Class() != "Function" {
			panic(fmt.Sprintf("field %s should be a function", name))
		}
	}
	*out = val
}

func (r *Runner) useDeployer1(call otto.FunctionCall) otto.Value {
	//TODO: implement deployer1
	return otto.NullValue()
}

func (r *Runner) useDeployer2(call otto.FunctionCall) otto.Value {
	//TODO: implement deployer2
	return otto.NullValue()
}

func (r *Runner) runScript(call otto.FunctionCall) otto.Value {
	shell := r.scriptShell
	if len(shell) == 0 {
		shell = []string{"/bin/bash"}
	}

	f := rg.Must(os.OpenFile(r.script, os.O_RDONLY, 0))
	defer f.Close()

	cmd := exec.Command(shell[0], shell[1:]...)
	cmd.Env = rg.Must(r.createEnviron())
	cmd.Stdin = f
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	rg.Must0(cmd.Run())

	return otto.NullValue()
}

func (r *Runner) runDockerBuild(call otto.FunctionCall) otto.Value {
	//TODO: implement package
	return otto.NullValue()
}

func (r *Runner) runDockerPush(call otto.FunctionCall) otto.Value {
	//TODO: implement publish
	return otto.NullValue()
}

func (r *Runner) useKubernetesWorkload(call otto.FunctionCall) otto.Value {
	if arg := call.Argument(0); arg.IsObject() {
		obj := arg.Object()
		r.loadObjectStringField(&r.workloadNamespace, obj, "namespace")
		r.loadObjectStringField(&r.workloadName, obj, "name")
		r.loadObjectStringField(&r.workloadKind, obj, "kind")
		r.loadObjectStringField(&r.workloadContainer, obj, "container")
		r.loadObjectBoolField(&r.workloadInit, obj, "init")
	}
	return rg.Must(r.createPlainObject(map[string]any{
		"namespace": r.workloadNamespace,
		"name":      r.workloadName,
		"kind":      r.workloadKind,
		"container": r.workloadContainer,
		"init":      r.workloadInit,
	}))
}

func (r *Runner) resolveCodingCredentials() (username string, password string) {
	var parts []string
	if r.codingValuesTeam != "" {
		parts = append(parts, SanitizeEnvName(r.codingValuesTeam))
		if r.codingValuesProject != "" {
			parts = append(parts, SanitizeEnvName(r.codingValuesProject))
			if r.codingValuesRepo != "" {
				parts = append(parts, SanitizeEnvName(r.codingValuesRepo))
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
		r.loadObjectStringField(&r.codingValuesTeam, obj, "team")
		r.loadObjectStringField(&r.codingValuesProject, obj, "project")
		r.loadObjectStringField(&r.codingValuesRepo, obj, "repo")
		r.loadObjectStringField(&r.codingValuesBranch, obj, "branch")
		r.loadObjectStringField(&r.codingValuesFile, obj, "file")
		r.loadObjectFunctionField(&r.codingValuesUpdate, obj, "update")
	}
	return rg.Must(r.createPlainObject(map[string]any{
		"team":    r.codingValuesTeam,
		"project": r.codingValuesProject,
		"repo":    r.codingValuesRepo,
		"branch":  r.codingValuesBranch,
		"file":    r.codingValuesFile,
		"update":  r.codingValuesUpdate,
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

	r.vm.Set("useEnv", r.createGetterSetterForObject(r.env, "env"))
	r.vm.Set("useDeployer1", r.useDeployer1)
	r.vm.Set("useDeployer2", r.useDeployer2)
	r.vm.Set("useRegistry", r.createGetterSetterForString(&r.registry, "registry"))
	r.vm.Set("useImage", r.createGetterSetterForString(&r.image, "image"))
	r.vm.Set("useProfile", r.createGetterSetterForString(&r.profile, "profile"))
	r.vm.Set("useVersion", r.createGetterSetterForString(&r.version, "version"))
	r.vm.Set("useDockerConfig", r.createGetterSetterForLongString(&r.dockerConfig, "docker config", func(buf []byte) (out string, err error) {
		_, out, err = r.createTempFile("config.json", bytes.TrimSpace(buf))
		return
	}))
	r.vm.Set("useKubeconfig", r.createGetterSetterForLongString(&r.kubeconfig, "kubeconfig", func(buf []byte) (out string, err error) {
		buf = bytes.TrimSpace(buf)
		if bytes.HasPrefix(buf, []byte("{")) {
			if buf, err = ConvertJSONToYAML(buf); err != nil {
				return
			}
		}
		out, _, err = r.createTempFile("kubeconfig.yaml", buf)
		return
	}))
	r.vm.Set("useScript", r.createGetterSetterForLongString(&r.script, "script", func(buf []byte) (out string, err error) {
		out, _, err = r.createTempFile("script.sh", bytes.TrimSpace(buf))
		return
	}))
	r.vm.Set("useScriptShell", r.createGetterSetterForStringSlice(&r.scriptShell, "script shell"))
	r.vm.Set("runScript", r.runScript)
	r.vm.Set("useDockerfile", r.createGetterSetterForString(&r.dockerfile, "dockerfile"))
	r.vm.Set("useDockerContext", r.createGetterSetterForString(&r.dockerContext, "docker context"))
	r.vm.Set("runDockerBuild", r.runDockerBuild)
	r.vm.Set("runDockerPush", r.runDockerPush)
	r.vm.Set("useKubernetesWorkload", r.useKubernetesWorkload)
	r.vm.Set("useCodingValues", r.useCodingValues)
	r.vm.Set("deployKubernetesWorkload", r.deployKubernetesWorkload)
	r.vm.Set("deployCodingValues", r.deployCodingValues)
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
	r.codingValuesUpdate = otto.Value{}
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
