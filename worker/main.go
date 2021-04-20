package worker

import (
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"time"
)

const (
	address = "192.168.0.11:50051"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Worker ", log.LstdFlags|log.Lshortfile)
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewBrokerClient(conn)

	broker, err := client.PinPong(context.Background())
	if err != nil {
		logger.Fatalln(err)
	}

	go func() {
		for {
			logger.Println("Wait for pong...")

			pong, err := broker.Recv()
			if err != nil {
				logger.Fatalln(err)
			}

			logger.Println("Pong", pong.Time)
		}
	}()

	for {
		err = broker.Send(&pb.Ping{
			Message: "Hallo",
			Time:    timestamppb.Now(),
		})

		if err != nil {
			logger.Fatalln(err)
		}

		time.Sleep(time.Hour * 20)
	}

	//ctx := context.Background()
	//rep, err := client.CreateWork(ctx, &pb.NewWorkRequest{
	//	SimulationName: "TicToc",
	//	Source:         nil,
	//	Configs:        nil,
	//})
	//
	//if err != nil {
	//	logger.Fatalf("could not greet: %v\n", err)
	//}
	//
	//logger.Printf("SimulationId: %s\n", rep.GetSimulationId())
	//
	//status, err := client.Status(ctx, &pb.StatusRequest{SimulationId: "abcd"})
	//if err != nil {
	//	logger.Fatalf("could not get status: %v", err)
	//}
	//
	//jsonBytes, _ := json.MarshalIndent(status, "", "    ")
	//logger.Printf("Status: %s\n", jsonBytes)
}
