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
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

type consumerId = string

type provider struct {
	pb.UnimplementedProviderServer
	providerId string
	store      *storage.Server

	mu          sync.RWMutex
	freeSlots   uint32
	requests    map[consumerId]uint32
	assignments map[consumerId]uint32
	allocate    map[consumerId]chan<- uint32
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

	log.Printf("Run: id=%v config='%s' runNum='%s'",
		simulation.Id, simulation.Config, simulation.RunNum)

	if simulation.Config == "" {
		err = fmt.Errorf("simulation config missing")
		return
	}

	if simulation.RunNum == "" {
		err = fmt.Errorf("simulation run number missing")
		return
	}

	if prov.freeSlots == 0 {
		err = fmt.Errorf("no free slots avilable")
		return
	}

	prov.mu.Lock()
	prov.freeSlots--
	prov.mu.Unlock()

	defer func() {
		prov.mu.Lock()
		prov.freeSlots++
		prov.mu.Unlock()
	}()

	return prov.run(simulation)
}

func (prov *provider) Allocate(stream pb.Provider_AllocateServer) (err error) {

	jobs, err := stream.Recv()
	if err != nil {
		log.Println(err)
		return
	}

	cId := fmt.Sprintf("%x", rand.Uint32())
	log.Printf("Allocate: register cId=%v", cId)

	allocate := make(chan uint32)
	prov.mu.Lock()
	prov.allocate[cId] = allocate
	prov.requests[cId] = jobs.Request
	prov.mu.Unlock()

	defer func() {
		log.Printf("Allocate: unregister cId=%v", cId)

		prov.mu.Lock()
		delete(prov.allocate, cId)
		delete(prov.requests, cId)
		prov.mu.Unlock()
		close(allocate)
	}()

	go func() {
		for slots := range allocate {
			log.Printf("Allocate: cId=%s allocate=%d", cId, slots)

			err = stream.Send(&pb.AllocatedSlots{
				Slots: slots,
			})
			if err != nil {
				log.Println(err)
				break
			}
		}
	}()

	for {
		jobs, err = stream.Recv()
		if err != nil {
			break
		}

		log.Printf("Allocate: cId=%s request=%d", cId, jobs.Request)

		prov.mu.Lock()
		prov.requests[cId] = jobs.Request
		prov.mu.Unlock()
	}

	return
}
