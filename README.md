# OMNeT++ simulation distributor

## Omnetpp

```
```

## Developer Notes

Install protobuf dependencies

```shell
go get -u google.golang.org/grpc
```

Generate protobuf

```shell
cd proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    broker.proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    storage.proto
```

## Worker

First configure your working client:

```shell
go run cmd/worker/worker.go --deviceName $(hostname -s) \
    --brokerAddress 192.168.0.11:50051 \
    --configure
```

Start a worker

```shell
go run cmd/worker/worker.go
```

## Simulation

Start a new simulation

```shell
go run cmd/simulation/simulation.go --path ~/Desktop/tictoc --configs TicToc18
```
