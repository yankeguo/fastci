useDockerImages(
    'yankeguo/debian:12'
)
useDockerContext('testdata/context1')
runDockerBuild()
runDockerPush()