package storage

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	storage pb.StorageClient
}

func (client *Client) Close() {
	_ = client.conn.Close()
}

func InitClient() (client Client) {
	conn, err := grpc.Dial(storageAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return
	}

	client.conn = conn
	client.storage = pb.NewStorageClient(conn)

	return
}
