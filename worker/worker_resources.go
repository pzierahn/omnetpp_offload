package worker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
)

func (client *workerConnection) SendWorkRequest(link pb.Broker_TaskSubscriptionClient) (err error) {

	info := pb.WorkRequest{
		WorkerId: client.workerId,
	}

	err = link.Send(&info)

	return
}
