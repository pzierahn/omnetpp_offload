package worker

import (
	pb "com.github.patrickz98.omnet/proto"
	"com.github.patrickz98.omnet/simple"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"runtime"
	"time"
)

func Link(config Config) (err error) {

	logger.Println("config", simple.PrettyString(config))

	//
	// Set up a connection to the server
	//

	conn, err := grpc.Dial(config.BrokerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewBrokerClient(conn)

	md := metadata.New(map[string]string{
		"workerId": config.WorkerId,
	})

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	link, err := client.Link(ctx)
	if err != nil {
		logger.Fatalln(err)
	}

	wClient := workerClient{
		config:        config,
		link:          link,
		freeResources: runtime.NumCPU(),
	}

	go func() {
		for {
			tasks, err := link.Recv()
			if err != nil {
				logger.Println(err)
				return
			}

			err = wClient.OccupyResource(len(tasks.Jobs))
			if err != nil {
				logger.Println(err)
				return
			}

			logger.Printf("received task %v\n", tasks.ProtoReflect())

			go runTasks(&wClient, tasks)
		}
	}()

	for {
		err = wClient.SendClientInfo()
		if err != nil {
			break
		}

		time.Sleep(time.Second * 23)
	}

	return
}
