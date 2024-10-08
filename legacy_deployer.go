package fastci

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/yankeguo/fastci/pkg/legacy_deployer_tmplfuncs"
	"github.com/yankeguo/rg"
	"gopkg.in/yaml.v3"
)

const (
	deployerProfileDefault = "default"
)

type deployerOptions struct {
	Cluster   string `json:"cluster"`
	Manifest  string `json:"manifest"`
	Profile   string `json:"profile"`
	Namespace string `json:"namespace"`
	Workload  string `json:"workload"`
	Image     string `json:"image"`
	Version   string `json:"version"`
}

func useDeployer(r *Runner, opts deployerOptions) (err error) {
	defer rg.Guard(&err)

	{
		jobNameSplits := strings.Split(strings.TrimSpace(os.Getenv("JOB_NAME")), ".")

		if len(jobNameSplits) >= 2 {
			if opts.Profile == "" {
				opts.Profile = jobNameSplits[len(jobNameSplits)-1]
			}
			if opts.Workload == "" {
				opts.Workload = jobNameSplits[len(jobNameSplits)-2]
			}
		}
		if len(jobNameSplits) >= 4 {
			if opts.Namespace == "" {
				opts.Namespace = jobNameSplits[len(jobNameSplits)-3]
			}
			if opts.Cluster == "" {
				opts.Cluster = jobNameSplits[len(jobNameSplits)-4]
			}
		}

		if opts.Image == "" {
			if opts.Namespace == "" {
				opts.Image = opts.Workload
			} else {
				opts.Image = opts.Namespace + "-" + opts.Workload
			}
		}

		if opts.Version == "" {
			opts.Version = os.Getenv("BUILD_NUMBER")
		}
		if opts.Version == "" {
			opts.Version = strconv.FormatInt(time.Now().Unix(), 10)
		}
	}

	if opts.Cluster == "" {
		err = errors.New("'cluster' is not set, either via JOB_NAME or useDeployer() options")
		return
	}

	if opts.Workload == "" {
		err = errors.New("'workload' is not set, either via JOB_NAME or useDeployer() options")
		return
	}

	if opts.Profile == "" {
		err = errors.New("'profile' is not set, either via JOB_NAME or useDeployer() options")
		return
	}

	var bufManifest []byte

	// try read the manifest file content
	if opts.Manifest == "" {
		if bufManifest, err = os.ReadFile("deployer.yml"); err != nil {
			if bufManifest, err = os.ReadFile("deployer.yaml"); err != nil {
				err = nil
			}
		}
	} else {
		bufManifest = rg.Must(os.ReadFile(opts.Manifest))
	}

	var header struct {
		Version int `yaml:"version"`
	}
	rg.Must0(yaml.Unmarshal(bufManifest, &header))

	if header.Version == 0 {
		err = useDeployerVersion1(r, opts, bufManifest)
	} else if header.Version == 2 {
		err = useDeployerVersion2(r, opts, bufManifest)
	} else {
		err = errors.New("unsupported deployer version")
	}
	return
}

func useDeployerVersion1(r *Runner, opts deployerOptions, manifest []byte) (err error) {
	defer rg.Guard(&err)

	// load runner state from preset
	{
		var bufPreset []byte

		home := rg.Must(os.UserHomeDir())

		if bufPreset, err = os.ReadFile(filepath.Join(home, ".deployer", "preset-"+opts.Cluster+".yml")); err != nil {
			if bufPreset, err = os.ReadFile(filepath.Join(home, ".deployer", "preset-"+opts.Cluster+".yaml")); err != nil {
				return
			}
		}

		var preset struct {
			Registry     string         `yaml:"registry"`
			Kubeconfig   map[string]any `yaml:"kubeconfig"`
			Dockerconfig map[string]any `yaml:"dockerconfig"`
		}

		rg.Must0(yaml.Unmarshal(bufPreset, &preset))

		if len(preset.Kubeconfig) > 0 {
			bufKubeconfig := rg.Must(yaml.Marshal(preset.Kubeconfig))
			r.state.kubernetes.kubeconfigPath, _ = rg.Must2(r.createTempFile("kubeconfig.yaml", bufKubeconfig))
			log.Printf("use kubeconfig: %s", r.state.kubernetes.kubeconfigPath)
		}

		if len(preset.Dockerconfig) > 0 {
			bufDockerConfig := rg.Must(json.Marshal(preset.Dockerconfig))
			_, r.state.docker.configPath = rg.Must2(r.createTempFile("config.json", bufDockerConfig))
			log.Printf("use docker config: %s", r.state.docker.configPath)
		}

		r.state.docker.images = append(r.state.docker.images, path.Join(preset.Registry, opts.Image+":"+opts.Profile))
		r.state.docker.images = append(r.state.docker.images, path.Join(preset.Registry, opts.Image+":"+opts.Profile+"-"+opts.Version))
	}

	// load runner state from manifest
	{
		type chunk struct {
			Build   []string          `yaml:"build"`
			Package []string          `yaml:"package"`
			Vars    map[string]string `yaml:"vars"`
		}

		var (
			root   chunk
			chunks map[string]chunk
		)

		rg.Must0(yaml.Unmarshal(manifest, &root))
		rg.Must0(yaml.Unmarshal(manifest, &chunks))

		var (
			finalBuild   []string
			finalPackage []string
			finalVars    = map[string]string{
				"profile": opts.Profile,
				"PROFILE": strings.ToUpper(opts.Profile),
			}
		)

		baseDir := "."

		if opts.Manifest != "" {
			baseDir = filepath.Base(opts.Manifest)
		}

		// compose build
		if len(chunks[opts.Profile].Build) > 0 {
			finalBuild = chunks[opts.Profile].Build
		} else if len(chunks[deployerProfileDefault].Build) > 0 {
			finalBuild = chunks[deployerProfileDefault].Build
		} else {
			finalBuild = root.Build
		}
		if len(finalBuild) == 0 {
			if buf, _ := os.ReadFile(filepath.Join(baseDir, "docker-build."+opts.Profile+".sh")); len(buf) > 0 {
				log.Println("use docker-build." + opts.Profile + ".sh")
				finalBuild = strings.Split(string(buf), "\n")
			}
		}

		// compose package
		if len(chunks[opts.Profile].Package) > 0 {
			finalPackage = chunks[opts.Profile].Package
		} else if len(chunks[deployerProfileDefault].Package) > 0 {
			finalPackage = chunks[deployerProfileDefault].Package
		} else {
			finalPackage = root.Package
		}
		if len(finalPackage) == 0 {
			if buf, _ := os.ReadFile(filepath.Join(baseDir, "Dockerfile."+opts.Profile)); len(buf) > 0 {
				log.Println("use Dockerfile." + opts.Profile)
				finalPackage = strings.Split(string(buf), "\n")
			}
		}

		// compose final vars
		for k, v := range root.Vars {
			finalVars[k] = v
			finalVars[k+"__uppercase"] = strings.ToUpper(v)
			finalVars[k+"__lowercase"] = strings.ToLower(v)
		}
		for k, v := range chunks[deployerProfileDefault].Vars {
			finalVars[k] = v
			finalVars[k+"__uppercase"] = strings.ToUpper(v)
			finalVars[k+"__lowercase"] = strings.ToLower(v)
		}
		for k, v := range chunks[opts.Profile].Vars {
			finalVars[k] = v
			finalVars[k+"__uppercase"] = strings.ToUpper(v)
			finalVars[k+"__lowercase"] = strings.ToLower(v)
		}

		// render build and package
		for k, v := range finalVars {
			src := fmt.Sprintf("{{__%s__}}", k)
			for i := range finalBuild {
				finalBuild[i] = strings.ReplaceAll(finalBuild[i], src, v)
			}
			for i := range finalPackage {
				finalPackage[i] = strings.ReplaceAll(finalPackage[i], src, v)
			}
		}

		if len(finalBuild) > 0 {
			content := strings.Join(finalBuild, "\n")
			file, _ := rg.Must2(r.createTempFile("script.sh", []byte(content)))
			log.Println("use build script:\n", content)
			r.state.script.path = file
		}

		if len(finalPackage) > 0 {
			content := strings.Join(finalPackage, "\n")
			file, _ := rg.Must2(r.createTempFile("Dockerfile", []byte(content)))
			log.Println("use Dockerfile:\n", content)
			r.state.docker.dockerfilePath = file
		}
	}

	return
}

func useDeployerVersion2(r *Runner, opts deployerOptions, manifest []byte) (err error) {
	defer rg.Guard(&err)

	// load runner state from preset
	{
		var bufPreset []byte

		home := rg.Must(os.UserHomeDir())

		if bufPreset, err = os.ReadFile(filepath.Join(home, ".deployer2", "preset-"+opts.Cluster+".yml")); err != nil {
			if bufPreset, err = os.ReadFile(filepath.Join(home, ".deployer2", "preset-"+opts.Cluster+".yaml")); err != nil {
				return
			}
		}

		var preset struct {
			Registry     string         `yaml:"registry"`
			Kubeconfig   map[string]any `yaml:"kubeconfig"`
			Dockerconfig map[string]any `yaml:"dockerconfig"`
		}

		rg.Must0(yaml.Unmarshal(bufPreset, &preset))

		if len(preset.Kubeconfig) > 0 {
			bufKubeconfig := rg.Must(yaml.Marshal(preset.Kubeconfig))
			r.state.kubernetes.kubeconfigPath, _ = rg.Must2(r.createTempFile("kubeconfig.yaml", bufKubeconfig))
			log.Printf("use kubeconfig: %s", r.state.kubernetes.kubeconfigPath)
		}

		if len(preset.Dockerconfig) > 0 {
			bufDockerConfig := rg.Must(json.Marshal(preset.Dockerconfig))
			_, r.state.docker.configPath = rg.Must2(r.createTempFile("config.json", bufDockerConfig))
			log.Printf("use docker config: %s", r.state.docker.configPath)
		}

		r.state.docker.images = append(r.state.docker.images, path.Join(preset.Registry, opts.Image+":"+opts.Profile))
		r.state.docker.images = append(r.state.docker.images, path.Join(preset.Registry, opts.Image+":"+opts.Profile+"-"+opts.Version))
	}

	// load runner state from manifest
	{
		type chunk struct {
			Build   []string          `yaml:"build"`
			Package []string          `yaml:"package"`
			Vars    map[string]string `yaml:"vars"`
		}

		var (
			chunks map[string]chunk
		)

		rg.Must0(yaml.Unmarshal(manifest, &chunks))

		var (
			finalBuild   []string
			finalPackage []string
			finalVars    = map[string]string{}
		)

		// compose build
		if len(chunks[opts.Profile].Build) > 0 {
			finalBuild = chunks[opts.Profile].Build
		} else if len(chunks[deployerProfileDefault].Build) > 0 {
			finalBuild = chunks[deployerProfileDefault].Build
		}

		// compose package
		if len(chunks[opts.Profile].Package) > 0 {
			finalPackage = chunks[opts.Profile].Package
		} else if len(chunks[deployerProfileDefault].Package) > 0 {
			finalPackage = chunks[deployerProfileDefault].Package
		}

		// compose final vars
		for k, v := range chunks[deployerProfileDefault].Vars {
			finalVars[k] = v
		}
		for k, v := range chunks[opts.Profile].Vars {
			finalVars[k] = v
		}

		data := map[string]interface{}{
			"Env":     rg.Must(r.createEnvironMap()),
			"Vars":    finalVars,
			"Profile": opts.Profile,
		}

		if len(finalBuild) > 0 {
			out := &bytes.Buffer{}
			rg.Must0(rg.Must(template.New("").Funcs(legacy_deployer_tmplfuncs.Funcs).Parse(strings.Join(finalBuild, "\n"))).Execute(out, data))
			content := out.String()

			file, _ := rg.Must2(r.createTempFile("script.sh", []byte(content)))
			log.Println("use build script:\n", content)
			r.state.script.path = file
		}

		if len(finalPackage) > 0 {
			out := &bytes.Buffer{}
			rg.Must0(rg.Must(template.New("").Funcs(legacy_deployer_tmplfuncs.Funcs).Parse(strings.Join(finalPackage, "\n"))).Execute(out, data))
			content := out.String()

			file, _ := rg.Must2(r.createTempFile("Dockerfile", []byte(content)))
			log.Println("use Dockerfile:\n", content)
			r.state.docker.dockerfilePath = file
		}
	}

	return
}
