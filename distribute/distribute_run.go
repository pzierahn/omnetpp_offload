package distribute

import (
	"context"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/provider"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"google.golang.org/grpc"
	"path/filepath"
)

func Run(gConf gconfig.GRPCConnection, config *Config) (err error) {

	if config.Tag == "" {
		config.Tag = filepath.Base(config.Path)
	}

	logger.Println("run new simulation")

	conn, err := grpc.Dial(gConf.DialAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	broker := pb.NewBrokerClient(conn)

	ctx := context.Background()
	simulationId, err := broker.Create(ctx, &pb.Simulation{
		Tag:       config.Tag,
		OppConfig: config.OppConfig,
	})

	if err != nil {
		return
	}

	logger.Printf("simulationId: %s", simulationId.Id)

	//
	// Clean source folder and upload source
	//

	compiler := provider.Compiler{
		Broker:         broker,
		Storage:        storage.ConnectClient(conn),
		SimulationId:   simulationId.Id,
		SimulationBase: config.Path,
		OppConfig:      config.OppConfig,
	}

	logger.Printf("cleaning: %s", config.Path)
	err = compiler.Clean()
	if err != nil {
		err = fmt.Errorf("couldn't clean simulation source: %v", err)
		return
	}

	ref, err := UploadSource(conn, simulationId.Id, config)
	if err != nil {
		err = fmt.Errorf("couldn't upload simulation source: %v", err)
		return
	}

	logger.Printf("uploaded to %v", ref)

	_, err = broker.SetSource(ctx, &pb.Source{
		SimulationId: simulationId.Id,
		Source:       ref,
	})
	if err != nil {
		err = fmt.Errorf("couldn't set simulation source: %v", err)
		return
	}

	//
	// Upload simulation source
	//

	err = compiler.Compile()
	if err != nil {
		err = fmt.Errorf("couldn't compile simulation: %v", err)
		return
	}

	err = compiler.CheckinBinary()
	if err != nil {
		err = fmt.Errorf("couldn't checkin simulation binary: %v", err)
		return
	}

	//
	// Extract configs and commit tasks
	//

	simConfs, err := extractConfigs(compiler.Opp())
	if err != nil {
		return
	}

	var items []*pb.SimulationRun

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

		for _, runNum := range runNums {
			items = append(items, &pb.SimulationRun{
				SimulationId: simulationId.Id,
				Config:       name,
				RunNumber:    runNum,
			})
		}
	}

	_, err = broker.AddTasks(ctx, &pb.Tasks{
		SimulationId: simulationId.Id,
		Items:        items,
	})
	if err != nil {
		err = fmt.Errorf("couldn't commit simulation tasks: %v", err)
		return
	}

	return
}
