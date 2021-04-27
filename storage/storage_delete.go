package storage

import (
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (client *Client) Delete(file *pb.StorageRef) (status *pb.StorageStatus, err error) {

	status, err = client.storage.Delete(context.Background(), file)
	return
}
