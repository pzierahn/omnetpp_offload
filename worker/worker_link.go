package worker

import (
	"context"
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/grpc/metadata"
	"runtime"
)

func (client *workerConnection) StartLink(ctx context.Context) (err error) {

	logger.Println("start worker", client.workerId)

	md := metadata.New(map[string]string{
		"workerId": client.workerId,
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"numCPU":   fmt.Sprint(runtime.NumCPU()),
	})

	ctx = metadata.NewOutgoingContext(ctx, md)

	// Link to the work stream
	link, err := client.broker.TaskSubscription(ctx)
	if err != nil {
		return
	}

	exit := make(chan bool)
	defer close(exit)

	for idx := 0; idx < client.agents; idx++ {

		//
		// Start worker agents
		//

		go func(idx int) {
			for {
				logger.Printf("agent %d waiting for work\n", idx)

				var task *pb.Task
				task, err = link.Recv()
				if err != nil {
					logger.Printf("agent %d: %v", idx, err)
					break
				}

				logger.Printf("agent %d received work (%s_%s_%s)\n",
					idx, task.SimulationId, task.Config, task.RunNumber)

				client.OccupyResource(1)
				client.runTasks(task)
				client.FeeResource()

				err = client.SendResourceCapacity(link)
				if err != nil {
					logger.Printf("agent %d: %v", idx, err)
					break
				}
			}

			logger.Printf("agent %d exiting\n", idx)
			exit <- true
		}(idx)
	}

	err = client.SendResourceCapacity(link)
	if err != nil {
		return
	}

	for idx := 0; idx < client.agents; idx++ {
		<-exit
	}

	logger.Println("closing connection to broker")

	err = client.Close()

	return
}
