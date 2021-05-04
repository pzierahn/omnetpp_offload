package broker

import (
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/utils"
	"google.golang.org/grpc/metadata"
)

func (server *broker) TaskSubscription(stream pb.Broker_TaskSubscriptionServer) (err error) {

	md, ok := metadata.FromIncomingContext(stream.Context())

	if !ok {
		logger.Printf("metadata missing")
		err = fmt.Errorf("metadata missing")
		return
	}

	var workerId string
	workerId, err = utils.MetaString(md, "workerId")
	if err != nil {
		return
	}

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
		var info *pb.ResourceCapacity
		info, err = stream.Recv()
		if err != nil {
			break
		}

		logger.Printf("%s freeResources=%v\n", workerId, info.FreeResources)

		server.db.SetCapacity(workerId, info)
		server.db.DistributeWork()
	}

	logger.Printf("lost connection to %s\n", workerId)

	return
}
