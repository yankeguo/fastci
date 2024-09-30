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

Alias `deployer1`

Use the `deployer1` preset and manifest, for compatibility with the legacy toolchain.

```javascript
useDeployer1("eco-staging", "deployer.yml");
```

#### `useDeployer2(preset, manifest="deployer.yml")`

Alias `deployer2`

Use the `deployer2` preset and manifest, for compatibility with the legacy toolchain.

```javascript
useDeployer2("eco-staging", "deployer2.yml");
```

#### `useRegistry(registry)`

Alias `registry`

Set the container registry for the pipeline.

```javascript
useRegistry("registry.cn-hangzhou.aliyuncs.com/eco-staging");
```

#### `useJenkins()`

Alias `jenkins`

Use the `Jenkins` environment variables for container image naming and environment switching.

```javascript
useJenkins();
```

#### `useDockerConfig(dockerConfig)`

Aliases `dockerConfig`, `useDockerconfig`, `dockerconfig`

Set the Docker configuration for the pipeline.

```javascript
useDockerConfig({ auths: {'registry.cn-hangzhou.aliyuncs.com': { username: "username", password: "password" } });
useDockerConfig('/path/to/.docker/dir');
```

#### `useKubeConfig(kubeConfig)`

Aliases `kubeConfig`, `useKubeconfig`, `kubeconfig`

Set the Kubernetes configuration for the pipeline.

```javascript
useKubeConfig({ apiVersion: "v1", clusters: [], contexts: [], users: [] });
useKubeConfig("/path/to/.kube/config/file");
```

## Credits

GUO YANKE, MIT License
