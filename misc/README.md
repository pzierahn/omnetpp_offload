# Build docker images

```
echo 'export PATH="$HOME/install/omnetpp/bin/:$PATH"' >> ~/.bashrc

docker build -t pzierahn/omnetpp:amd64 .
docker build -t pzierahn/omnetpp:arm64 .

docker push pzierahn/omnetpp:amd64
docker push pzierahn/omnetpp:arm64

docker buildx build \
    --push \
    --platform linux/arm64,linux/amd64 \
    --tag pzierahn/omnetpp:latest -f Dockerfile.omnetpp .
```
