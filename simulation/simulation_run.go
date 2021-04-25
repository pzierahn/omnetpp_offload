package simulation

import (
	"context"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/defines"
	"github.com/patrickz98/project.go.omnetpp/omnetpp"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"google.golang.org/grpc"
	"sort"
)

// Connect to broker to commit a new simulation
func commitSimulation(simulation *pb.Simulation) (err error) {

	conn, err := grpc.Dial(defines.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewBrokerClient(conn)

	ctx := context.Background()
	reply, err := client.ExecuteSimulation(ctx, simulation)
	if err != nil {
		return
	}

	logger.Println("simulationId", reply.SimulationId)

	return
}

func Run(config *Config) (err error) {

	if config.SimulationId == "" {
		err = fmt.Errorf("no SimulationId in config")
		return
	}

	//
	// Extract configurations and clean the project afterwards
	//

	omnet := omnetpp.New(&config.Config)

	err = omnet.Setup()
	if err != nil {
		err = fmt.Errorf("couldn't setup simulation: %v", err)
		return
	}

	simConfs, err := extractConfigs(omnet)
	if err != nil {
		return
	}

	var confs []*pb.Simulation_RunConfig

	if len(config.SimulateConfigs) == 0 {
		// Add all configs
		for conf := range simConfs {
			config.SimulateConfigs = append(config.SimulateConfigs, conf)
		}
	}

	for _, name := range config.SimulateConfigs {

		runNums, ok := simConfs[name]
		if !ok {
			err = fmt.Errorf("unknown simulation configuration '%s'\n", name)
			return
		}

		confs = append(confs, &pb.Simulation_RunConfig{
			Config:     name,
			RunNumbers: runNums,
		})
	}

	// Sort OCD
	sort.Slice(confs, func(i, j int) bool {
		return fmt.Sprintf("%s-%-3s", confs[i].Config, confs[i].RunNumbers) >
			fmt.Sprintf("%s-%-3s", confs[j].Config, confs[j].RunNumbers)
	})

	//
	// Clean and upload simulation
	//

	err = omnet.Clean()
	if err != nil {
		err = fmt.Errorf("couldn't clean simulation: %v", err)
		return
	}

	ref, err := Upload(config)
	if err != nil {
		err = fmt.Errorf("couldn't upload simulation: %v", err)
		return
	}

	simulation := pb.Simulation{
		SimulationId: config.SimulationId,
		Tag:          config.Tag,
		OppConfig:    config.OppConfig,
		Source:       ref,
		Run:          confs,
	}

	simple.WritePretty("debug/simulation.json", &simulation)

	err = commitSimulation(&simulation)

	return
}
