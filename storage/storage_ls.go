package storage

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
)

func (client *Client) List(file *pb.StorageRef) (list *pb.StorageList, err error) {
	list, err = client.storage.List(context.Background(), file)
	return
}
