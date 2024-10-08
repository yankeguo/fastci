useDockerImages(
    'yankeguo/debian:12'
)
useDockerBuildContext('testdata/context1')
runDockerBuild()
runDockerPush()