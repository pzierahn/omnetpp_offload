# OMNeT++ simulation distributor

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    broker.proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    storage.proto

scp -r project.go.omnet ubuntu@raspberry3b:~/
```

## Add a new worker

```
go run cmd/worker.go -deviceName $(hostname -s) \
    -brokerAddress 192.168.0.11:50051 \
    -configure
    

```
