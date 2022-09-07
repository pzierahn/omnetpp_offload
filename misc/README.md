# Build OMNeT++ Docker image

The official OMNeT++ Docker images don't support arm64 architectures.
Hence, an arm64 supported image needs to be build manually.

The Dockerfile.omnetpp is copied and compiled from
the [official OMNeT++ Dockerfiles](https://github.com/omnetpp/dockerfiles).

```shell
docker buildx build \
    --push \
    --platform linux/arm64,linux/amd64 \
    --tag pzierahn/omnetpp:6.0.1 \
    --file Dockerfile.omnetpp .
```
