useDockerImages(
    'yankeguo/debian:12'
)
useDockerfile(
    'FROM debian:12'
)
useDockerBuildContext('testdata/context1')
runDockerBuild()
runDockerPush()