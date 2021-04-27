package distribute

import (
	"context"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	"github.com/patrickz98/project.go.omnetpp/omnetpp"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/grpc"
	"sort"
)

// Connect to broker to commit a new simulation
func commitSimulation(config gconfig.GRPCConnection, simulation *pb.Simulation) (err error) {

	conn, err := grpc.Dial(config.DialAddr(), grpc.WithInsecure(), grpc.WithBlock())
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

func Run(conn gconfig.GRPCConnection, config *Config) (err error) {

	config.generateId()

	logger.Println("simulationId", config.SimulationId)

	//
	// Extract configurations and clean the project afterwards
	//

	// Todo: change config.Config.Config
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

	//simple.WritePretty("debug/simulation.json", &simulation)

	err = commitSimulation(conn, &simulation)

	return
}
