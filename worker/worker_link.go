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

	// channel for thread safe error communication
	done := make(chan bool)
	defer close(done)

	go func() {

		//
		// Receive task from the broker
		//

		for {
			var tasks *pb.Tasks
			tasks, err = link.Recv()
			if err != nil {
				break
			}

			client.OccupyResource(len(tasks.Items))

			//logger.Printf("received task %v\n", tasks)

			for _, job := range tasks.Items {
				go func(job *pb.Task) {
					logger.Printf("running task %v_%v_%v\n", job.SimulationId, job.Config, job.RunNumber)
					client.runTasks(job)

					logger.Printf("free resource %v_%v_%v\n", job.SimulationId, job.Config, job.RunNumber)
					client.FeeResource()

					err = client.SendResourceCapacity(link)
					if err != nil {
						logger.Fatalln(err)
					}
				}(job)
			}
		}

		done <- true
	}()

	go func() {

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

		done <- true
	}()

	<-done

	return
}
