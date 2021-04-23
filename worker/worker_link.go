package worker

import (
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"runtime"
	"time"
)

func (client *workerConnection) StartLink(ctx context.Context) (err error) {

	md := metadata.New(map[string]string{
		"workerId": client.config.WorkerId,
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"numCPU":   fmt.Sprint(runtime.NumCPU()),
	})

	ctx = metadata.NewOutgoingContext(ctx, md)

	// Link to the work stream
	link, err := client.client.Link(ctx)
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

			client.OccupyResource(len(tasks.Jobs))

			//logger.Printf("received task %v\n", tasks)

			for _, job := range tasks.Jobs {
				go func(job *pb.Work) {
					logger.Printf("running task %v_%v\n", job.Config, job.RunNumber)
					client.runTasks(job)

					logger.Printf("free resource %v_%v\n", job.Config, job.RunNumber)

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
		// This will prevent the connection from closing the connection
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
