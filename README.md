# OMNeT++ simulation distributor

## Install tools

```shell
rm -rf ~/go/bin/opp_edge_*

go install cmd/broker/opp_edge_broker.go
go install cmd/config/opp_edge_config.go
go install cmd/distribute/opp_edge_run.go
go install cmd/omnetpp/opp_edge_opp.go
go install cmd/storage/opp_edge_storage.go
go install cmd/worker/opp_edge_worker.go
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

Start a worker

```shell
go run cmd/worker/worker.go
```

## Simulation

Start a new simulation

```shell
go run cmd/simulation/simulation.go --path ~/Desktop/tictoc --configs TicToc18
```


## Example simulations

```
go run cmd/distribute/opp_edge_run.go -path ../TaskletSimulator -config ../TaskletSimulator/opp-edge-config.json
go run cmd/distribute/opp_edge_run.go -path ~/Desktop/tictoc -config ~/Desktop/tictoc/opp-edge-config.json
```

```
GOOS=linux GOARCH=amd64 go build cmd/ice/ice.go
```