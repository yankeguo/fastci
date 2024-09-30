# fastci

An intuitive CLI tool that encompasses the entire cycle of `build`, `package`, `push`, and `deploy` operations.

## Usage

`fastci` is desiged to read the `JavaScript` pipeline from `stdin`, and execute in a single process.

```shell
cat <<-EOF | fastci
useDeployer2("eco-staging", "mobile/deployer2.yml")
useJenkins()
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
```

#### `useDeployer2(preset, manifest="deployer.yml")`

Use the `deployer2` preset and manifest, for compatibility with the legacy toolchain.

```javascript
useDeployer2("eco-staging", "mobile/deployer2.yml");
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

### `useVersion(version)`

Get or set the container version for the pipeline.

```javascript
useVersion("1.0.0");
```

#### `useJenkins()`

Use the `Jenkins` environment variables for `image`, `profile` and `version`.

```javascript
useJenkins();
```

#### `useDockerConfig(dockerConfig)`

Set the Docker configuration for the pipeline.

```javascript
useDockerConfig({
  auths: {
    "registry.cn-hangzhou.aliyuncs.com": {
      username: "username",
      password: "password",
    },
  },
});
useDockerConfig("/path/to/docker/config/dir");
```

#### `useKubeconfig(kubeconfig)`

Set the Kubernetes configuration for the pipeline.

```javascript
useKubeconfig({ apiVersion: "v1", clusters: [], contexts: [], users: [] });
useKubeconfig("/path/to/kubeconfig/file");
```

#### `useBuildScript(script)`

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

#### `useBuildScriptFile(scriptFile)`

Get or set the build script file for the pipeline.

```javascript
useBuildScriptFile("build.sh");
```

### `useBuildScriptShell(shell)`

Get or set the build script shell for the pipeline.

```javascript
useBuildScriptShell("zsh");
```

#### `useBuildDir(buildDir)`

Get or set the build directory for the pipeline.

```javascript
useBuildDir("build");
```

#### `doBuild()`

Execute the build script in the build directory.

```javascript
doBuild();
```

## Credits

GUO YANKE, MIT License
