package provider

import (
	"bytes"
	"context"
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
		return
	}

	path := filepath.Join(cachePath, bundle.SimulationId)

	err = simple.UnTarGz(cachePath, bytes.NewReader(byt))
	if err != nil {
		_ = os.RemoveAll(path)
	}

	return
}

func (prov *provider) Compile(_ context.Context, simulation *pb.Simulation) (bin *pb.Binary, err error) {
	log.Printf("Compile: %v", simulation.Id)
	return prov.compile(simulation)
}

func (prov *provider) Run(_ context.Context, run *pb.SimulationRun) (ref *pb.StorageRef, err error) {
	return
}
