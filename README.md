# fastci

An intuitive CLI tool that encompasses the entire cycle of `build`, `package`, `push`, and `deploy` operations.

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

### Variables

#### `env`

The environment variables that will be used in the pipeline.

```javascript
console.log(env["CI_COMMIT_REF_NAME"]);
```

### Functions

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

#### `useRegistry(registry)`

Get or set the container registry for the pipeline.

```javascript
useRegistry("registry.cn-hangzhou.aliyuncs.com/eco-staging");
```

#### `useImage(image)`

Get or set the container image name

```javascript
useImage("my-project");
```

#### `useProfile(profile)`

Get or set the container profile for the pipeline.

```javascript
useProfile("staging");
```

#### `useVersion(version)`

Get or set the container version for the pipeline.

```javascript
useVersion("1.0.0");
```

#### `useDockerConfig(dockerConfig)`

Set the Docker configuration for the pipeline.

```javascript
useDockerConfig({
  content: {
    auths: {
      "registry.cn-hangzhou.aliyuncs.com": {
        username: "username",
        password: "password",
      },
    },
  },
  // content: '',
  // base64: '',
  // path: '',
});
useDockerConfig({ path: "/path/to/docker/config/dir" });
```

#### `useKubeconfig(kubeconfig)`

Set the Kubernetes configuration for the pipeline.

```javascript
useKubeconfig({
  content: { apiVersion: "v1", clusters: [], contexts: [], users: [] },
  // content: '',
  // base64: '',
  // path: '',
});
useKubeconfig({ path: "/path/to/kubeconfig/file" });
```

#### `useBuildScript(buildScript)`

Get or set the build script for the pipeline.

```javascript
// script content
useBuildScript(
  "\
#!/bin/bash\n\
echo 'Building...'\n\
",
);

// script object
useBuildScript({
  // content: '',
  // base64: '',
  // path: '',
});
```

#### `useBuildScriptShell(shell)`

Get or set the shell for the build script.

```javascript
useBuildScriptShell("zsh");
```

#### `doBuild()`

Execute the previous configured script in the pipeline.

```javascript
doBuild();
```

#### `usePackageDockerfile(dockerfile)`

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

## Credits

GUO YANKE, MIT License
