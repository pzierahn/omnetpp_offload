package simulation

import (
	"com.github.patrickz98.omnet/defines"
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"google.golang.org/grpc"
)

func Run(config Config) {

	ref, err := Upload(config)

	conn, err := grpc.Dial(defines.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewBrokerClient(conn)

	simulation := pb.Simulation{
		SimulationId: config.Id,
		Source:       ref,
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
