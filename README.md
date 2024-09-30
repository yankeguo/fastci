# fastci

An intuitive CLI tool that encompasses the entire cycle of `build`, `package`, `push`, and `deploy` operations.

## Usage

`fastci` is desiged to read the `JavaScript` pipeline from `stdin`, and execute in a single process.

```shell
cat <<-EOF | fastci
useDeployer2("eco-staging", "mobile/deployer2.yml")
runScript()
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

#### `useScript(script)`

Get or set the build script for the pipeline.

```javascript
// script content
useScript(
  "\
#!/bin/bash\n\
echo 'Building...'\n\
",
);

// script object
useScript({
  // content: '',
  // base64: '',
  // path: '',
});
```

#### `useScriptShell(shell)`

Get or set the shell for the build script.

```javascript
useScriptShell("zsh");
```

#### `runScript()`

Execute the previous configured script in the pipeline.

```javascript
runScript();
```

## Credits

GUO YANKE, MIT License
