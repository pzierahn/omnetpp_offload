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
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type simulationId = string

type provider struct {
	pb.UnimplementedProviderServer
	providerId string
	store      *storage.Server

	cond        *sync.Cond
	slots       uint32
	freeSlots   int32
	requests    map[simulationId]uint32
	assignments map[simulationId]uint32
	runCtx      map[simulationId]context.CancelFunc
	allocate    map[simulationId]chan<- uint32
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

	if run.SimulationId == "" {
		err = fmt.Errorf("SimulationId missing")
		return
	}

	if atomic.LoadInt32(&prov.freeSlots) == 0 {
		err = fmt.Errorf("no free slots avilable")
		return
	}

	// Simulation Run Id = srId
	srId := fmt.Sprintf("%s_%s_%s", run.SimulationId, run.Config, run.RunNum)

	ctx, cnl := context.WithTimeout(context.Background(), time.Minute*60)
	defer cnl()

	cond := prov.cond

	cond.L.Lock()
	prov.runCtx[srId] = cnl
	cond.L.Unlock()

	atomic.AddInt32(&prov.freeSlots, -1)

	defer func() {
		log.Printf("Run: id=%v config='%s' runNum='%s' done",
			run.SimulationId, run.Config, run.RunNum)

		cond.L.Lock()
		atomic.AddInt32(&prov.freeSlots, 1)
		prov.assignments[run.SimulationId]--
		delete(prov.runCtx, srId)
		cond.Broadcast()
		cond.L.Unlock()
	}()

	return prov.run(ctx, run)
}

func (prov *provider) Allocate(stream pb.Provider_AllocateServer) (err error) {

	var sId string

	allocate := make(chan uint32)

	defer func() {
		log.Printf("Allocate: unregister SimulationId=%v", sId)

		cond := prov.cond

		cond.L.Lock()
		delete(prov.allocate, sId)
		delete(prov.requests, sId)
		delete(prov.assignments, sId)

		// Cancel running simulations
		for id, cnl := range prov.runCtx {
			if strings.HasPrefix(id, sId) {
				log.Printf("Allocate: cancel %s", id)
				cnl()
			}
		}

		// Clean up and remove simulation (delete simulation bucket)
		_, _ = prov.store.Drop(nil, &pb.BucketRef{Bucket: sId})

		cond.Broadcast()
		cond.L.Unlock()

		close(allocate)
	}()

	go func() {
		for slots := range allocate {
			log.Printf("Allocate: SimulationId=%s allocate=%d", sId, slots)

			err = stream.Send(&pb.AllocatedSlots{Slots: slots})
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

		if sId == "" {
			sId = req.SimulationId
			log.Printf("Allocate: register SimulationId=%v", sId)

			cond.L.Lock()
			prov.allocate[sId] = allocate
			cond.L.Unlock()
		}

		if sId == "" {
			err = fmt.Errorf("error: missing SimulationId")
			log.Println(err)
			return
		}

		log.Printf("Allocate: SimulationId=%s request=%d", sId, req.Request)

		cond.L.Lock()

		if val, _ := prov.requests[sId]; val != req.Request {
			prov.requests[sId] = req.Request
			cond.Broadcast()
		}

		cond.L.Unlock()
	}

	return
}
