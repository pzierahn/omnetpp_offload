package provider

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"github.com/pzierahn/project.go.omnetpp/sysinfo"
	"log"
	"os"
	"path/filepath"
)

type provider struct {
	pb.UnimplementedProviderServer
	providerId string
	store      *storage.Server
}

func (prov *provider) Info(_ context.Context, _ *pb.Empty) (info *pb.ProviderInfo, err error) {

	log.Printf("Info:")
	info = prov.info()

	return
}

func (prov *provider) Status(_ context.Context, _ *pb.Empty) (util *pb.Utilization, err error) {

	log.Printf("Status:")

	util, err = sysinfo.GetUtilization()
	return
}

func (prov *provider) Checkout(_ context.Context, bundle *pb.Bundle) (empty *pb.Empty, err error) {

	log.Printf("Checkout: %v", bundle.SimulationId)

	empty = &pb.Empty{}

	byt, err := prov.store.Get(bundle.Source)
	if err != nil {
		log.Printf("Checkout: error %v", err)
		return
	}

	path := filepath.Join(cachePath, bundle.SimulationId)

	err = simple.UnTarGz(cachePath, bytes.NewReader(byt))
	if err != nil {
		log.Printf("Checkout: error %v", err)
		_ = os.RemoveAll(path)
	}

	return
}

func (prov *provider) Compile(_ context.Context, simulation *pb.Simulation) (bin *pb.Binary, err error) {
	log.Printf("Compile: %v", simulation.Id)
	return prov.compile(simulation)
}

func (prov *provider) ListRunNums(_ context.Context, simulation *pb.Simulation) (runs *pb.SimulationRuns, err error) {

	log.Printf("ListRunNums: id=%v config='%s'", simulation.Id, simulation.Config)

	if simulation.Config == "" {
		err = fmt.Errorf("simulation config missing")
		return
	}

	_, opp := newOpp(simulation)

	runNums, err := opp.GetRunNumbers(simulation.Config)
	if err != nil {
		return
	}

	runs = &pb.SimulationRuns{
		SimulationId: simulation.Id,
		Config:       simulation.Config,
		Runs:         runNums,
	}

	return
}

func (prov *provider) Run(_ context.Context, simulation *pb.Simulation) (ref *pb.StorageRef, err error) {

	log.Printf("ListRunNums: id=%v config='%s' runNum='%s'",
		simulation.Id, simulation.Config, simulation.RunNum)

	if simulation.Config == "" {
		err = fmt.Errorf("simulation config missing")
		return
	}

	if simulation.RunNum == "" {
		err = fmt.Errorf("simulation run number missing")
		return
	}

	return prov.run(simulation)
}
