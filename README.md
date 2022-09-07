# OMNeT++ simulation distributor

## Install command line tools

```shell
go install cmd/worker/opp_offload_worker.go
go install cmd/run/opp_offload_run.go
go install cmd/config/opp_offload_config.go
go install cmd/broker/opp_offload_broker.go
```

## Install and run with Docker

```shell
docker pull pzierahn/omnetpp_offload
docker run --rm pzierahn/omnetpp_offload opp_offload_worker -broker 85.214.35.83 -name `hostname -s`
```

## Build and upload docker images

Build cross-platform images for amd64 and arm64.

```shell
docker buildx build \
    --push \
    --platform linux/arm64,linux/amd64 \
    --tag pzierahn/omnetpp_offload:latest .
```

> Build alternative: ```docker build -t pzierahn/omnetpp_offload .```

## Run example simulations

```shell
go run cmd/run/opp_offload_run.go -path ~/github/TaskletSimulator
go run cmd/run/opp_offload_run.go -path evaluation/tictoc
```

## Install and run broker

```shell
go install cmd/broker/opp_offload_broker.go
nohup opp_offload_broker > opp_offload_broker.log 2>&1 &
```

## Developer Notes

Install protobuf dependencies.

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

go get -u google.golang.org/grpc
GOOS=linux GOARCH=amd64 go build cmd/consumer/opp_edge_run.go
```
