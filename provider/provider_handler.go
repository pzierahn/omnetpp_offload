package provider

import (
	"context"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"github.com/patrickz98/project.go.omnetpp/sysinfo"
	"log"
)

type provider struct {
	pb.UnimplementedProviderServer
	providerId string
	config     gconfig.Worker
	storage    storage.Client
	agents     int
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

func (prov *provider) Init(_ context.Context, simulation *pb.Simulation) (res *pb.Empty, err error) {
	return
}

func (prov *provider) CompileSync(_ context.Context, simulation *pb.SimulationId) (bin *pb.Binary, err error) {
	return
}

func (prov *provider) RunSync(_ context.Context, run *pb.SimulationRun) (ref *pb.StorageRef, err error) {
	return
}
