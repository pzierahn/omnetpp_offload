# OMNeT++ simulation distributor

## Install and run a worker

```
docker pull pzierahn/omnetpp_edge
docker run --rm pzierahn/omnetpp_edge opp_edge_worker -broker 85.214.35.83 -name `hostname -s`
```

## Build and upload docker images

Build cross-platform images for amd and arm

```shell
docker buildx build \
    --push \
    --platform linux/arm64,linux/amd64 \
    --tag pzierahn/omnetpp_edge:latest .

docker pull pzierahn/omnetpp_edge
docker run --rm pzierahn/omnetpp_edge opp_edge_worker -broker 31.18.129.212 -name `hostname -s`
```

> Build alternative: ```docker build -t pzierahn/omnetpp_edge .```

## Run example simulations

```
go run cmd/consumer/opp_edge_run.go -path ~/github/TaskletSimulator -config ~/github/TaskletSimulator/opp-edge-config.json
go run cmd/consumer/opp_edge_run.go -path ~/Desktop/tictoc -config ~/Desktop/tictoc/opp-edge-config.json
```

## Developer Notes

Install protobuf dependencies

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

go get -u google.golang.org/grpc
GOOS=linux GOARCH=amd64 go build cmd/consumer/opp_edge_run.go
```

Generate protobufs

```shell
cd proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    broker.proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    storage.proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    opp_config.proto
```
