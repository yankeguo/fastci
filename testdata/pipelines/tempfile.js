useDockerConfig({
    auths: {
        "https://index.docker.io/v1/": {}
    }
})

if (!dockerconfig()) {
    throw new Error('dockerconfig() failed')
}

useKubeconfig({ hello: 'world' })

if (!kubeconfig()) {
    throw new Error('kubeconfig() failed')
}