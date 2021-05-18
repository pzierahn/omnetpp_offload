package worker

import (
	"context"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/sysinfo"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"runtime"
	"time"
)

func (client *workerConnection) StartLink(ctx context.Context) (err error) {

	logger.Println("start worker", client.providerId)

	md := NewDeviceInfo(client.providerId).MarshallMeta()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.broker.WorkSubscription(ctx)
	if err != nil {
		return
	}

	workStream := make(chan *pb.Work)

	for inx := 0; inx < client.agents; inx++ {
		go func(inx int) {
			for work := range workStream {
				switch task := work.Work.(type) {
				case *pb.Work_Task:
					logger.Printf("(%d) Run simulation %v", inx, task)
					client.runTasks(task)

				case *pb.Work_Compile:
					logger.Printf("(%d) Compile simulation %v", inx, task)
					client.compile(task)
				}
			}
		}(inx)
	}

	go func() {
		var work *pb.Work

		for {
			logger.Printf("Waiting for work...")

			work, err = stream.Recv()
			if err != nil {
				panic(err)
			}

			logger.Printf("Received work %v", work)
			workStream <- work
		}

	}()

	for range time.Tick(time.Second) {
		usage := sysinfo.GetCPUUsage()

		// logger.Printf("Sending usage=%f", usage)

		state := &pb.ProviderState{
			ProviderId:  client.providerId,
			CpuUsage:    float32(usage),
			MemoryUsage: 0,
			Arch: &pb.OsArch{
				Os:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
			NumCPUs: uint32(runtime.NumCPU()),
			Tasks:   nil,
			Compile: "",
			Updated: timestamppb.Now(),
		}

		err = stream.Send(state)
		if err != nil {
			return
		}
	}

	logger.Println("closing connection to broker")

	return
}
