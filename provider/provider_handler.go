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
	"sync"
	"sync/atomic"
)

type consumerId = string

type provider struct {
	pb.UnimplementedProviderServer
	providerId string
	store      *storage.Server

	cond        *sync.Cond
	slots       uint32
	freeSlots   int32
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

func (prov *provider) Run(_ context.Context, run *pb.SimulationRun) (ref *pb.StorageRef, err error) {

	log.Printf("Run: id=%v config='%s' runNum='%s'",
		run.SimulationId, run.Config, run.RunNum)

	if run.Config == "" {
		err = fmt.Errorf("simulation config missing")
		return
	}

	if run.RunNum == "" {
		err = fmt.Errorf("simulation run number missing")
		return
	}

	if run.ConsumerId == "" {
		err = fmt.Errorf("ConsumerId missing")
		return
	}

	if atomic.LoadInt32(&prov.freeSlots) == 0 {
		err = fmt.Errorf("no free slots avilable")
		return
	}

	cond := prov.cond

	atomic.AddInt32(&prov.freeSlots, -1)

	defer func() {
		log.Printf("Run: id=%v config='%s' runNum='%s' done",
			run.SimulationId, run.Config, run.RunNum)

		cond.L.Lock()
		atomic.AddInt32(&prov.freeSlots, 1)
		prov.assignments[run.ConsumerId]--
		cond.Broadcast()
		cond.L.Unlock()
	}()

	return prov.run(run)
}

func (prov *provider) Allocate(stream pb.Provider_AllocateServer) (err error) {

	var cId string

	allocate := make(chan uint32)

	defer func() {
		log.Printf("Allocate: unregister ConsumerId=%v", cId)

		cond := prov.cond

		cond.L.Lock()
		delete(prov.allocate, cId)
		delete(prov.requests, cId)
		delete(prov.assignments, cId)

		cond.Broadcast()
		cond.L.Unlock()

		close(allocate)

		// TODO: clean up and remove simulation
	}()

	go func() {
		for slots := range allocate {
			log.Printf("Allocate: ConsumerId=%s allocate=%d", cId, slots)

			err = stream.Send(&pb.AllocatedSlots{
				Slots: slots,
			})
			if err != nil {
				log.Println(err)
				break
			}
		}
	}()

	cond := prov.cond

	for {
		var req *pb.AllocateRequest
		req, err = stream.Recv()
		if err != nil {
			break
		}

		if cId == "" {
			cId = req.ConsumerId
			log.Printf("Allocate: register ConsumerId=%v", cId)

			cond.L.Lock()
			prov.allocate[cId] = allocate
			cond.L.Unlock()
		}

		if cId == "" {
			err = fmt.Errorf("error: missing ConsumerId")
			log.Println(err)
			return
		}

		log.Printf("Allocate: ConsumerId=%s request=%d", cId, req.Request)

		cond.L.Lock()

		if val, _ := prov.requests[cId]; val != req.Request {
			prov.requests[cId] = req.Request
			cond.Broadcast()
		}

		cond.L.Unlock()
	}

	return
}
