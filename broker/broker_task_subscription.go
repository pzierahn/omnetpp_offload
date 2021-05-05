package broker

import (
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/worker"
	"google.golang.org/grpc/metadata"
)

func (server *broker) TaskSubscription(stream pb.Broker_TaskSubscriptionServer) (err error) {

	md, ok := metadata.FromIncomingContext(stream.Context())

	if !ok {
		logger.Printf("metadata missing")
		err = fmt.Errorf("metadata missing")
		return
	}

	var workerInfo worker.DeviceInfo
	workerInfo.UnMarshallMeta(md)

	workerId := workerInfo.WorkerId
	if workerId == "" {
		err = fmt.Errorf("missing workerId in metadata")
		return
	}

	logger.Printf("connected worker %s (%v, %v, %v)",
		workerId, workerInfo.Os, workerInfo.Arch, workerInfo.NumCPUs)

	//
	// Send work to clients
	//

	workStream := server.db.NewWorker(workerId)
	defer server.db.RemoveWorker(workerId)

	go func() {
		for {
			job, ok := <-workStream
			if !ok {
				logger.Printf("exit work subscription for %s\n", workerId)
				break
			}

			logger.Printf("send %v work to %s", taskId(job), workerId)

			err := stream.Send(job)
			if err != nil {
				logger.Println(err)
				break
			}
		}
	}()

	//
	// Receive Client Info
	//

	for {
		var workReq *pb.WorkRequest
		workReq, err = stream.Recv()
		if err != nil {
			break
		}

		logger.Printf("work request from %s\n", workReq.WorkerId)

		server.db.IncreaseCapacity(workReq.WorkerId)
		server.db.DistributeWork()
	}

	logger.Printf("lost connection to %s\n", workerId)

	return
}
