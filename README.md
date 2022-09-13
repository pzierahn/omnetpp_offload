# OMNeT++ simulation distributor

## Command line tools and usage

### Install

```shell
go install cmd/worker/opp_offload_worker.go
go install cmd/run/opp_offload_run.go
go install cmd/config/opp_offload_config.go
go install cmd/broker/opp_offload_broker.go
```

### opp_offload_config

Helps you to create a global and local omnetpp_offload configuration files

```
Usage of opp_offload_config:
  -broker string
        set broker address
  -jobs int
        set how many jobs should be started (default 8)
  -name string
        set worker name (default "Patricks-MBP")
  -paths
        print paths
  -port int
        set broker port (default 8888)
  -save
        persist config globally
  -stargate int
        set stargate port (default 8889)
```

### opp_offload_broker

The broker connects providers and workers

```
Usage of opp_offload_broker:
  -broker string
        set broker address
  -port int
        set broker port (default 8888)
  -stargate int
        set stargate port (default 8889)
```

### opp_offload_worker

The worker starts an OMNeT++ work provider for you

```
Usage of opp_offload_worker:
  -broker string
        set broker address
  -clean
        clean all cache files
  -jobs int
        set how many jobs should be started (default 8)
  -name string
        set worker name (default "Patricks-MBP")
  -port int
        set broker port (default 8888)
  -stargate int
        set stargate port (default 8889)
```

### opp_offload_run

opp_offload_run will offload simulations for you

```
Usage of opp_offload_run:
  -broker string
        set broker address
  -config string
        set simulation config JSON
  -path string
        set simulation path (default ".")
  -port int
        set broker port (default 8888)
  -stargate int
        set stargate port (default 8889)
  -timeout duration
        set timeout for execution (default 3h0m0s)
```

### Install and run with Docker

* Get Docker image ```docker pull pzierahn/omnetpp_offload```
* Start worker ```docker run --rm pzierahn/omnetpp_offload opp_offload_worker -broker 85.214.35.83 -name `hostname -s```

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
GOOS=linux GOARCH=amd64 go build cmd/consumer/opp_offload_run.go
```
