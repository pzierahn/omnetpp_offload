package provider

import (
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/sysinfo"
	"google.golang.org/grpc/metadata"
	"time"
)

func (client *workerConnection) StartLink(ctx context.Context) (err error) {

	logger.Println("start worker", client.providerId)

	md := NewDeviceInfo(client.providerId).MarshallMeta()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.broker.Assignments(ctx)
	if err != nil {
		return
	}

	assignments := make(chan *pb.Assignment)

	for inx := 0; inx < client.agents; inx++ {
		go func(inx int) {
			for assignment := range assignments {
				switch task := assignment.Do.(type) {
				case *pb.Assignment_Run:
					logger.Printf("(%d) Run simulation %v", inx, task)
					//client.runTasks(task)

				case *pb.Assignment_Build:
					logger.Printf("(%d) Compile simulation %v", inx, task)
					//client.compile(task)
				}
			}
		}(inx)
	}

	go func() {
		var work *pb.Assignment

		for {
			//logger.Printf("Waiting for work...")

			work, err = stream.Recv()
			if err != nil {
				panic(err)
			}

			//logger.Printf("Received work %v", work)
			assignments <- work
		}

	}()

	for range time.Tick(time.Second * 1) {
		var usage *pb.Utilization
		usage, err = sysinfo.GetUtilization()

		// logger.Printf("Sending usage=%v", usage)

		err = stream.Send(usage)
		if err != nil {
			return
		}
	}

	logger.Println("closing connection to broker")

	return
}
