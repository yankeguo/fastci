# fastci

An intuitive CLI tool that encompasses the entire cycle of `build`, `package`, `publish`, and `deploy` operations.

[![Go](https://github.com/yankeguo/fastci/actions/workflows/go.yml/badge.svg)](https://github.com/yankeguo/fastci/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/yankeguo/fastci/graph/badge.svg?token=91hTz3G4x3)](https://codecov.io/gh/yankeguo/fastci)

## Usage

`fastci` is designed to read the `JavaScript` pipeline from `stdin`, and execute in a single process.

```shell
cat <<-EOF | fastci
useDeployer2("eco-staging", "mobile/deployer2.yml")
doBuild()
EOF
```

## Pipeline

### General Conventions

#### Function with Long Text

For functions with long text, `fastci` supports `plain text`, `plain object`, `array of lines`, `base64`, and `file path`.

For example, the `useKubeconfig` function can be used in the following ways:

```javascript
useKubeconfig(
  "\
apiVersion: v1\n\
kind: aaaa\n\
",
);

useKubeconfig({
  content:
    "\
apiVersion: v1\n\
kind: aaaa\n\
",
});

useKubeconfig({
  content: {
    apiVersion: "v1",
    kind: "xxx",
  },
});

useKubeconfig({
  base64: "xxx",
});

useKubeconfig(["apiVersion: v1", "kind: xxx"]);

useKubeconfig("apiVersion: v1", "kind: xxx");
```

### Basic Functions

#### `useEnv`

Get or set the environment variables for the pipeline.

```javascript
// get all environment variables
useEnv();
// get the value of the specified environment variable
useEnv("key");
// set the value of the specified environment variable
useEnv("key", "value");
```

#### `useRegistry(registry)`

Get or set the container registry for the pipeline.

```javascript
// set the container registry
useRegistry("registry.cn-hangzhou.aliyuncs.com/eco-staging");
// get the container registry
useRegistry();
```

#### `useImage(image)`

Get or set the container base image name

```javascript
// set the container base image name
useImage("my-project");
// get the container base image name
useImage();
```

#### `useProfile(profile)`

Get or set the profile for the pipeline.

```javascript
// set the profile
useProfile("staging");
// get the profile
useProfile();
```

#### `useVersion(version)`

Get or set the container version for the pipeline.

```javascript
// set the container version
useVersion("114514");
// get the container version
useVersion();
// set the version to environment variable BUILD_NUMBER
useVersion(useEnv("BUILD_NUMBER"));
```

#### `useDockerConfig(dockerConfig)`

This function is **Long Text Supported**

Set the Docker configuration for the pipeline.

```javascript
// set the Docker config to the specified content
useDockerConfig({
  // content can be an object or a string
  content: {
    auths: {
      "registry.cn-hangzhou.aliyuncs.com": {
        username: "username",
        password: "password",
      },
    },
  },
});

// get the Docker config directory
useDockerConfig();
```

#### `useKubeconfig(kubeconfig)`

This function is **Long Text Supported**

Set the Kubernetes configuration for the pipeline.

```javascript
// set the Kubernetes config to the specified content
useKubeconfig({
  // content can be an object or a string
  content: {
    apiVersion: "v1",
    clusters: [],
    contexts: [],
    users: [],
  },
});

// get the Kubernetes config file path
useKubeconfig();
```

### Build Functions

#### `useBuildScript(buildScript)`

This function is **Long Text Supported**

Get or set the build script for the pipeline.

```javascript
// script content
useBuildScript(
  "\
#!/bin/bash\n\
echo 'Building...'\n\
",
);
```

#### `useBuildShell(shell)`

Get or set the shell for the build script.

```javascript
useBuildShell("zsh");
```

#### `doBuild()`

Execute the previous configured build script in the pipeline.

```javascript
doBuild();
```

### Package Functions

#### `usePackageDockerfile(dockerfile)`

This function is **Long Text Supported**

Get or set the Dockerfile for the package operation.

```javascript
usePackageDockerfile([
  "FROM ubuntu:24.04",
  'RUN echo "Hello, World!"',
  'CMD ["echo", "Hello, World!"]',
]);
```

#### `usePackageContext(context)`

Get or set the context for the package operation.

```javascript
usePackageContext("./docker/");
```

#### `doPackage()`

Package the container image with the previous configured Dockerfile.

```javascript
doPackage();
```

### Publish Functions

#### `doPublish()`

Push the container image to the registry.

```javascript
doPublish();
```

### Deploy Functions

#### `useKubernetesWorkload(opts)`

Configure the Kubernetes workload for deploying the container image.

```javascript
useKubernetesWorkload({
  // namespace of the workload
  namespace: "my-ns",
  // kind of the workload
  kind: "Deployment",
  // name of the workload
  name: "my-app",
  // if not set, container name will be the same as the workload name
  container: "my-app",
  // if it's a init container
  init: false,
});

// sub-sequence calls will merge the options
// for example, the following two ways are equivalent
useKubernetesWorkload({
  namespace: "my-ns",
  kind: "Deployment",
});
useKubernetesWorkload({
  name: "my-app",
});
//
//
useKubernetesWorkload({
  namespace: "my-ns",
  kind: "Deployment",
  name: "my-app",
});
```

#### `doDeployKubernetesWorkload()`

Deploy the container image to the Kubernetes cluster.

```javascript
doDeployKubernetesWorkload();
```

#### `useCodingValuesFile(opts)`

Configure the `coding.net` repository values file for deploying the container image.

```javascript
useCodingValuesFile({
  team: "my-team",
  project: "my-project",
  repo: "my-repo",
  branch: "main",
  file: "values.yaml",
  update: function (m) {
    m[useEnv("JOB_NAME")] = useEnv("BUILD_NUMBER");
  },
});

// sub-sequence calls will merge the options
// for example, the following two ways are equivalent
useCodingValuesFile({
  team: "my-team",
  project: "my-project",
  branch: "main",
  file: "values.yaml",
  update: function (m) {
    m[useEnv("JOB_NAME")] = useEnv("BUILD_NUMBER");
  },
});
useCodingValuesFile({
  repo: "my-repo",
});
//
useCodingValuesFile({
  team: "my-team",
  project: "my-project",
  repo: "my-repo",
  branch: "main",
  file: "values.yaml",
  update: function (m) {
    m[useEnv("JOB_NAME")] = useEnv("BUILD_NUMBER");
  },
});
```

#### `doDeployCodingValuesFile(options)`

Deploy the container image to patch `coding.net` repository values file.

```javascript
doDeployCodingValuesFile();
```

### Compatibility Functions

#### `useDeployer1(preset, manifest="deployer.yml")`

Use the `deployer1` preset and manifest, for compatibility with the legacy toolchain.

```javascript
useDeployer1("eco-staging", "deployer.yml");
useDeployer1("eco-staging");
useDeployer1({
  preset: "eco-staging",
  manifest: "deployer.yml",
});
```

#### `useDeployer2(preset, manifest="deployer.yml")`

Use the `deployer2` preset and manifest, for compatibility with the legacy toolchain.

```javascript
useDeployer2("eco-staging", "mobile/deployer2.yml");
useDeployer2("eco-staging");
useDeployer2({
  preset: "eco-staging",
  manifest: "deployer.yml",
});
```

## Credits

GUO YANKE, MIT License
