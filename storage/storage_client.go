package storage

import (
	"github.com/patrickz98/project.go.omnetpp/gconfig"
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

func ConnectClient(conn *grpc.ClientConn) (client Client) {
	client.conn = conn
	client.storage = pb.NewStorageClient(conn)

	return
}

func InitClient(server gconfig.GRPCConnection) (client Client) {
	conn, err := grpc.Dial(server.DialAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return
	}

	client.conn = conn
	client.storage = pb.NewStorageClient(conn)

	return
}
