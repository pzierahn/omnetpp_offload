package provider

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/sysinfo"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
)

type simulationId = string

func (prov *provider) GetSession(ctx context.Context, sim *pb.Simulation) (sess *pb.Session, err error) {

	log.Printf("GetSession: %v", sim.Id)

	prov.mu.Lock()
	defer prov.mu.Unlock()

	var ok bool
	if sess, ok = prov.sessions[sim.Id]; ok {
		return
	}

	sess = &pb.Session{
		SimulationId: sim.Id,
	}

	if deadline, ok := ctx.Deadline(); ok {
		sess.Ttl = timestamppb.New(deadline)
		go prov.expireSession(sim.Id, deadline)
	}

	prov.sessions[sim.Id] = sess
	prov.persistSessions()

	return
}

func (prov *provider) SetSession(_ context.Context, sess *pb.Session) (*pb.Session, error) {

	log.Printf("SetSession: %v", sess.SimulationId)

	prov.mu.Lock()
	prov.sessions[sess.SimulationId] = sess
	prov.persistSessions()
	prov.mu.Unlock()

	return sess, nil
}

func (prov *provider) Info(_ context.Context, _ *pb.Empty) (info *pb.ProviderInfo, err error) {

	log.Printf("Info:")
	info = prov.info()

	return
}

func (prov *provider) Status(ctx context.Context, _ *pb.Empty) (util *pb.Utilization, err error) {

	log.Printf("Status:")
	util, err = sysinfo.GetUtilization(ctx)
	return
}

func (prov *provider) Extract(_ context.Context, bundle *pb.Bundle) (empty *pb.Empty, err error) {

	log.Printf("Extract: %v %v", bundle.SimulationId, bundle.Source.Filename)

	empty = &pb.Empty{}

	byt, err := prov.store.Get(bundle.Source)
	if err != nil {
		log.Printf("Extract: %v error %v", bundle.SimulationId, err)
		return
	}

	path := filepath.Join(cachePath, bundle.SimulationId)

	err = simple.ExtractTarGz(cachePath, bytes.NewReader(byt))
	if err != nil {
		log.Printf("Extract: %v error %v", bundle.SimulationId, err)
		_ = os.RemoveAll(path)
	}

	return
}

func (prov *provider) Compile(ctx context.Context, simulation *pb.Simulation) (bin *pb.Binary, err error) {
	log.Printf("Compile: %v", simulation.Id)
	return prov.compile(ctx, simulation)
}

func (prov *provider) ListRunNums(ctx context.Context, simulation *pb.Simulation) (runs *pb.SimulationRuns, err error) {

	log.Printf("ListRunNums: id=%v config='%s'", simulation.Id, simulation.Config)

	if simulation.Config == "" {
		err = fmt.Errorf("simulation config missing")
		return
	}

	_, opp := newOpp(simulation)

	runNums, err := opp.GetRunNumbers(ctx, simulation.Config)
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

func (prov *provider) Run(ctx context.Context, run *pb.SimulationRun) (ref *pb.StorageRef, err error) {

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

	atomic.AddInt32(&prov.freeSlots, -1)

	defer func() {
		log.Printf("Run: id=%v config='%s' runNum='%s' done",
			run.SimulationId, run.Config, run.RunNum)

		cond := prov.cond
		cond.L.Lock()
		atomic.AddInt32(&prov.freeSlots, 1)

		if val, ok := prov.assignments[run.SimulationId]; ok && val > 0 {
			prov.assignments[run.SimulationId]--
		}

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

			prov.mu.Lock()
			prov.allocate[sId] = allocate
			prov.mu.Unlock()
		}

		if sId == "" {
			err = fmt.Errorf("error: missing SimulationId")
			log.Println(err)
			break
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
