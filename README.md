# OMNeT++ simulation distributor

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    broker.proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    storage.proto

scp -r project.go.omnetpp ubuntu@raspberry3b:~/

go get -u google.golang.org/grpc
```

## Worker

First configure your working client:
```
go run cmd/worker.go --deviceName $(hostname -s) \
    --brokerAddress 192.168.0.11:50051 \
    --configure
```


## Simulation
Options:
```
  -configs string
    	simulation config names
  -debug
    	send debug request
  -name string
    	name of the simulation
  -path string
    	path to OMNeT++ simulation
```
