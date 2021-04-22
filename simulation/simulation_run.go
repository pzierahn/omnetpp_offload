package simulation

import (
	"com.github.patrickz98.omnet/defines"
	"com.github.patrickz98.omnet/omnetpp"
	pb "com.github.patrickz98.omnet/proto"
	"com.github.patrickz98.omnet/simple"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"sort"
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

	allConfs, err := extractConfigs(omnet)
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

	if len(config.Configs) == 0 {

		//
		// Add all configs
		//

		for conf := range allConfs {
			config.Configs = append(config.Configs, conf)
		}
	}

	var confs []*pb.Config

	for _, name := range config.Configs {

		runNums, ok := allConfs[name]
		if !ok {
			fmt.Printf("unknown simulation configuration '%s'\n", name)
			os.Exit(1)
		}

		confs = append(confs, &pb.Config{
			Name:       name,
			RunNumbers: runNums,
		})
	}

	// Sort OCD
	sort.Slice(confs, func(i, j int) bool {
		return confs[i].Name > confs[j].Name
	})

	simulation := pb.Simulation{
		SimulationId: config.Id,
		Source:       ref,
		Configs:      confs,
	}

	simple.WritePretty("debug/simulation.json", &simulation)

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
