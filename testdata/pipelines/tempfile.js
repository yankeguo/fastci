useDockerConfig({
    auths: {
        "https://index.docker.io/v1/": {}
    }
})

if (!useDockerConfig()) {
    throw new Error('dockerconfig() failed')
}

useKubeconfig({ hello: 'world' })

if (!useKubeconfig()) {
    throw new Error('kubeconfig() failed')
}