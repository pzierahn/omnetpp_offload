package worker

import (
	"context"
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/grpc/metadata"
	"runtime"
	"time"
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
					break
				}

				logger.Printf("agent %d received work (%s_%s_%s)\n",
					idx, task.SimulationId, task.Config, task.RunNumber)

				client.OccupyResource(1)
				client.runTasks(task)
				client.FeeResource()

				err = client.SendResourceCapacity(link)
				if err != nil {
					logger.Fatalln(err)
				}
			}
		}(idx)
	}

	//
	// Send every 23 seconds the resource capacity
	// This will prevent the connection from closing
	//

	for {
		err = client.SendResourceCapacity(link)
		if err != nil {
			break
		}

		time.Sleep(time.Second * 23)
	}

	return
}
