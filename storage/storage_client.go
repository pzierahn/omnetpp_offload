package storage

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"google.golang.org/grpc"
)

type Client struct {
	storage pb.StorageClient
}

func FromClient(storeClient pb.StorageClient) (client Client) {
	client.storage = storeClient
	return
}

func FromConnection(conn *grpc.ClientConn) (client Client) {
	client.storage = pb.NewStorageClient(conn)
	return
}
