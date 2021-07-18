package storage

import (
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
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

func InitClient() (client Client) {
	conn, err := grpc.Dial(
		gconfig.BrokerDialAddr(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return
	}

	client.storage = pb.NewStorageClient(conn)

	return
}
