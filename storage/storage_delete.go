package storage

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
)

func (client *Client) Delete(file *pb.StorageRef) (status *pb.StorageStatus, err error) {

	status, err = client.storage.Delete(context.Background(), file)
	return
}
