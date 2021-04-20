package simulation

import (
	"com.github.patrickz98.omnet/defines"
	"com.github.patrickz98.omnet/omnetpp"
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func Run(config Config) {

	//
	// Extract configurations and clean the project afterwards
	//

	omnet := omnetpp.New(config.Path)

	err := omnet.Setup()
	if err != nil {
		err = fmt.Errorf("couldn't setup simulation: %v", err)
		return
	}

	confs, err := extractConfigs(omnet)
	if err != nil {
		return
	}

	err = omnet.Clean()
	if err != nil {
		err = fmt.Errorf("couldn't clean simulation: %v", err)
		return
	}

	ref, err := Upload(config)
	if err != nil {
		logger.Fatalf("couldn't upload simulation: %v", err)
	}

	//
	// Connect to broker to commit a new simulation
	//

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
		Configs:      confs,
	}

	ctx := context.Background()
	reply, err := client.NewSimulation(ctx, &simulation)
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("simulationId", reply.SimulationId)
}
