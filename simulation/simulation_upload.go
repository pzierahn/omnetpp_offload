package simulation

import (
	"com.github.patrickz98.omnet/defines"
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Simulation ", log.LstdFlags|log.Lshortfile)
}

func Run(filepath string) {

	conn, err := grpc.Dial(defines.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewBrokerClient(conn)

	simulation := pb.Simulation{
		SimulationId: "tictoc-1234",
		Source: &pb.StorageRef{
			Bucket:   "tictoc",
			Filename: "source.tar.gz",
		},
		Configs: []*pb.Config{
			{
				Name: "TicToc18",
				RunNumbers: []string{
					"1", "2", "3",
				},
			},
		},
	}

	ctx := context.Background()
	reply, err := client.NewSimulation(ctx, &simulation)
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("simulationId", reply.SimulationId)
}
