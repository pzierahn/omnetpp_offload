package worker

import (
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/grpc/metadata"
	"sync"
)

func (client *workerConnection) StartLink(ctx context.Context) (err error) {

	logger.Println("start worker", client.workerId)

	md := NewDeviceInfo(client.workerId).MarshallMeta()
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Link to the work stream
	link, err := client.broker.TaskSubscription(ctx)
	if err != nil {
		return
	}
	defer func() { _ = link.CloseSend() }()

	exit := make(chan bool)
	defer close(exit)

	work := make(chan *pb.Task)
	go func() {

		//
		// Single thread to receive tasks
		//

		for {
			task, err := link.Recv()
			if err != nil {
				logger.Printf("work receiver: %v", err)
				break
			}

			logger.Printf("receive work %v_%v_%v", task.SimulationId, task.Config, task.RunNumber)
			work <- task
		}

		logger.Printf("exit work receiver")
		close(work)
	}()

	var sendMu sync.Mutex

	for idx := 0; idx < client.agents; idx++ {

		//
		// Start worker agents
		//

		go func(idx int) {
			for {
				logger.Printf("agent %d send work request", idx)

				sendMu.Lock()
				err = client.SendWorkRequest(link)
				if err != nil {
					logger.Printf("agent %d: %v", idx, err)
					break
				}
				sendMu.Unlock()

				logger.Printf("agent %d waiting for work", idx)
				task, ok := <-work
				if !ok {
					break
				}

				logger.Printf("agent %d received work (%s_%s_%s)",
					idx, task.SimulationId, task.Config, task.RunNumber)

				client.runTasks(task)
			}

			logger.Printf("agent %d exiting", idx)
			exit <- true
		}(idx)
	}

	for idx := 0; idx < client.agents; idx++ {
		<-exit
	}

	logger.Println("closing connection to broker")

	return
}
