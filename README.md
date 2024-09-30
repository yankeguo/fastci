# fastci

An intuitive CLI tool that encompasses the entire cycle of `build`, `package`, `push`, and `deploy` operations.

## Usage

**Stages**

- `build`: Build the project.
- `package`: Package the project to container image.
- `push`: Push the container image to the registry.
- `deploy`: Deploy the container image, directly to the Kubernetes cluster or by updating GitOps repository.

**CLI**

Since CI pipeline is very complex, the `fastci` is desiged to read the `JavaScript` pipeline from stdin, and execute in a single process.

```shell
cat <<-EOF | fastci
use_deployer2_preset("eco-staging")
use_deployer2_manifest()
use_jenkins()
build()
package()
push()
deploy_to_k8s()
EOF
```

## Pipeline

### Variables

#### `env`

The environment variables that will be used in the pipeline.

```javascript
console.log(env["CI_COMMIT_REF_NAME"]);
```

### Configuration Functions

#### `useDeployer1(preset, manifest)`, `deployer1`

Use the `deployer1` preset and manifest, for compatibility with the legacy toolchain.

```javascript
useDeployer1("eco-staging", "deployer.yml");
```

#### `useDeployer2(preset, manifest)`, `deployer2`

Use the `deployer2` preset and manifest, for compatibility with the legacy toolchain.

```javascript
useDeployer2("eco-staging", "deployer2.yml");
```

#### `useJenkins()`, `jenkins`

Use the `Jenkins` environment variables for container image naming and environment switching.

```javascript
useJenkins();
```

## Credits

GUO YANKE, MIT License
