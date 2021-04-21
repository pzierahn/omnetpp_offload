package simulation

import (
	"com.github.patrickz98.omnet/defines"
	"com.github.patrickz98.omnet/omnetpp"
	pb "com.github.patrickz98.omnet/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func commitSimulation(simulation *pb.Simulation) (err error) {
	//
	// Connect to broker to commit a new simulation
	//

	conn, err := grpc.Dial(defines.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewBrokerClient(conn)

	ctx := context.Background()
	reply, err := client.NewSimulation(ctx, simulation)
	if err != nil {
		return
	}

	logger.Println("simulationId", reply.SimulationId)

	return
}

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
	// Create and push simulation to broker
	//

	simulation := pb.Simulation{
		SimulationId: config.Id,
		Source:       ref,
		Configs:      confs,
	}

	if err = commitSimulation(&simulation); err != nil {
		logger.Fatalln(err)
	}
}

func DebugRequest() {

	var mockConfigs []*pb.Config

	for inx := 0; inx < 200; inx++ {
		mockConfigs = append(mockConfigs, &pb.Config{
			Name: fmt.Sprintf("DebugConfig-%d", inx),
			RunNumbers: []string{
				"1",
				"2",
			},
		})
	}

	simulation := pb.Simulation{
		SimulationId: "debug-123456",
		Source: &pb.StorageRef{
			Bucket:   "abcd",
			Filename: "source",
		},
		Configs: mockConfigs,
	}

	if err := commitSimulation(&simulation); err != nil {
		logger.Fatalln(err)
	}
}
