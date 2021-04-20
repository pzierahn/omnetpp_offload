package worker

import (
	pb "com.github.patrickz98.omnet/proto"
	"com.github.patrickz98.omnet/simple"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	go func() {
		for {
			work, err := link.Recv()
			if err != nil {
				logger.Println(err)
				return
			}

			byt, err := json.MarshalIndent(work, "", "    ")
			if err != nil {
				logger.Println(err)
				return
			}

			logger.Println("work: ", string(byt))
		}
	}()

	for {
		logger.Println("sending info")

		info := pb.ClientInfo{
			Id:            config.WorkerId,
			Os:            runtime.GOOS,
			Arch:          runtime.GOARCH,
			NumCPU:        int32(runtime.NumCPU()),
			Timestamp:     timestamppb.Now(),
			FreeResources: 0,
		}

		err = link.Send(&info)

		time.Sleep(time.Second * 23)
	}
}
