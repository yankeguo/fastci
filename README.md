# fastci

An intuitive CLI tool that encompasses the entire cycle of `build`, `package`, `publish`, and `deploy` operations.

[![Go](https://github.com/yankeguo/fastci/actions/workflows/go.yml/badge.svg)](https://github.com/yankeguo/fastci/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/yankeguo/fastci/graph/badge.svg?token=91hTz3G4x3)](https://codecov.io/gh/yankeguo/fastci)

## Usage

`fastci` is designed to read the `JavaScript` pipeline from `stdin`, and execute in a single process.

```shell
cat <<-EOF | fastci
useDeployer2("eco-staging", "mobile/deployer2.yml")
runScript()
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

### Basic Configuration

#### `useShell(shell...)`

Get or set the shell for the current script.

```javascript
useShell("zsh");
useShell("bash", "-eux");
useShell(["bash", "-eux"]);
```

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

### Script

#### `useScript(script)`

This function is **Long Text Supported**

Get or set the current script for the pipeline.

```javascript
// script content
useScript(
  "\
#!/bin/bash\n\
echo 'Building...'\n\
",
);
```

#### `runScript()`

Execute the current script in the pipeline.

```javascript
runScript();
```

### Docker Build

#### `useDockerImages(images...)`

Get or set the Docker images for `docker build` and `docker push`

```javascript
useDockerImages("my-custom/ubuntu:24.04", "my-custom/ubuntu:24");
useDockerImages(["my-custom/ubuntu:24.04", "my-custom/ubuntu:24"]);
// clear the Docker images
useDockerImages(null);
useDockerImages([]);
```

#### `useDockerfile(dockerfile)`

This function is **Long Text Supported**

Get or set the current Dockerfile.

```javascript
useDockerfile([
  "FROM ubuntu:24.04",
  'RUN echo "Hello, World!"',
  'CMD ["echo", "Hello, World!"]',
]);
```

#### `useDockerContext(context)`

Get or set the current docker context.

```javascript
useDockerContext("./docker/");
```

#### `runDockerBuild()`

Package the container image with docker build command.

```javascript
runDockerBuild();
```

Returns the container images as array of string.

### Docker Push

#### `runDockerPush()`

Push the container image to the registry.

```javascript
runDockerPush();
```

### Deploy to Kubernetes

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

#### `deployKubernetesWorkload()`

Deploy the container image to the Kubernetes cluster.

```javascript
deployKubernetesWorkload();
```

### Deploy to Coding Values file

#### `useCodingValues(opts)`

Configure the `coding.net` repository values file for deploying the container image.

```javascript
useCodingValues({
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
useCodingValues({
  team: "my-team",
  project: "my-project",
  branch: "main",
  file: "values.yaml",
  update: function (m) {
    m[useEnv("JOB_NAME")] = useEnv("BUILD_NUMBER");
  },
});
useCodingValues({
  repo: "my-repo",
});
//
useCodingValues({
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

#### `deployCodingValues(options)`

Deploy the container image to patch `coding.net` repository values file.

```javascript
deployCodingValues();
```

`fastci` will search `coding.net` credentials from the environment variables, in order:

- `CODING_MY_TEAM_MY_PROJECT_MY_REPO_USERNAME` and `CODING_MY_TEAM_MY_PROJECT_MY_REPO_PASSWORD`
- `CODING_MY_TEAM_MY_PROJECT_USERNAME` and `CODING_MY_TEAM_MY_PROJECT_PASSWORD`
- `CODING_MY_TEAM_USERNAME` and `CODING_MY_TEAM_PASSWORD`
- `CODING_USERNAME` and `CODING_PASSWORD`

All characters in the team, project, and repo names will be converted to uppercase.

All non-numeric and non-alphabetic characters in the team, project, and repo names will be replaced with `_`.

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
