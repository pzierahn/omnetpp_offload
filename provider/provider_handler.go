package provider

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/sysinfo"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func (prov *provider) GetSession(ctx context.Context, sim *pb.Simulation) (sess *pb.Session, err error) {

	log.Printf("GetSession: %v", sim.Id)

	prov.mu.Lock()
	defer prov.mu.Unlock()

	var ok bool
	if sess, ok = prov.sessions[sim.Id]; !ok {
		//
		// Create new session
		//

		sess = &pb.Session{
			SimulationId: sim.Id,
			OppConfig:    sim.OppConfig,
		}

		if deadline, ok := ctx.Deadline(); ok {
			sess.Ttl = timestamppb.New(deadline)
			go prov.expireSession(sim.Id, deadline)
		}

		prov.sessions[sim.Id] = sess
		prov.persistSessions()
	}

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

func (prov *provider) Info(_ context.Context, _ *emptypb.Empty) (info *pb.ProviderInfo, err error) {

	log.Printf("Info:")
	info = prov.info()

	return
}

func (prov *provider) Status(ctx context.Context, _ *emptypb.Empty) (util *pb.Utilization, err error) {

	log.Printf("Status:")
	util, err = sysinfo.GetUtilization(ctx)
	return
}

func (prov *provider) Extract(_ context.Context, bundle *pb.Bundle) (empty *emptypb.Empty, err error) {

	log.Printf("Extract: %v %v", bundle.SimulationId, bundle.Source.Filename)

	empty = &emptypb.Empty{}

	byt, err := prov.store.PullFile(bundle.Source)
	if err != nil {
		log.Printf("Extract: %v error %v", bundle.SimulationId, err)
		return
	}

	path := filepath.Join(cachePath, bundle.SimulationId)

	err = simple.ExtractTarGz(cachePath, byt)
	if err != nil {
		log.Printf("Extract: %v error %v", bundle.SimulationId, err)
		_ = os.RemoveAll(path)
	}

	return
}

func (prov *provider) Compile(ctx context.Context, simulation *pb.Simulation) (bin *pb.Binary, err error) {
	log.Printf("Compile: %v", simulation.Id)
	started := time.Now()

	defer func() {
		duration := time.Now().Sub(started)
		prov.mu.Lock()
		prov.executionTimes[simulation.Id] += duration
		prov.mu.Unlock()
	}()

	return prov.compile(ctx, simulation)
}

func (prov *provider) ListRunNums(ctx context.Context, simulation *pb.Simulation) (runs *pb.SimulationRunList, err error) {

	log.Printf("ListRunNums: id=%v config='%s'", simulation.Id, simulation.Config)

	if simulation.Config == "" {
		err = fmt.Errorf("simulation config missing")
		return
	}

	_, opp := newOpp(simulation)

	runNums, err := opp.QRunNumbers(ctx, simulation.Config)
	if err != nil {
		return
	}

	runs = &pb.SimulationRunList{}

	for _, runNum := range runNums {
		runs.Items = append(runs.Items, &pb.SimulationRun{
			SimulationId: simulation.Id,
			Config:       simulation.Config,
			RunNum:       runNum,
		})
	}

	return
}

func (prov *provider) Run(ctx context.Context, run *pb.SimulationRun) (ref *pb.StorageRef, err error) {

	log.Printf("Run: id=%v config=%s runNum=%s",
		run.SimulationId, run.Config, run.RunNum)

	if run.SimulationId == "" {
		err = fmt.Errorf("simulation id missing")
		return
	}

	if run.Config == "" {
		err = fmt.Errorf("simulation config missing")
		return
	}

	if run.RunNum == "" {
		err = fmt.Errorf("simulation run number missing")
		return
	}

	started := time.Now()

	defer func() {
		log.Printf("Run: id=%v config=%s runNum=%s finished",
			run.SimulationId, run.Config, run.RunNum)

		duration := time.Now().Sub(started)
		prov.mu.Lock()
		prov.executionTimes[run.SimulationId] += duration
		prov.mu.Unlock()
	}()

	return prov.run(ctx, run)
}

func (prov *provider) Allocate(stream pb.Provider_AllocateServer) (err error) {

	ctx := stream.Context()
	simId, err := simple.MetaStringFromContext(ctx, "simulationId")
	if err != nil {
		log.Println(err)
		return
	}

	allocRecv := make(chan int, 1)
	defer close(allocRecv)

	prov.register(simId, allocRecv)
	defer prov.unregister(simId)

	var mu sync.Mutex
	var allocations uint32
	defer func() {
		mu.Lock()
		defer mu.Unlock()

		log.Printf("Allocate: feed back %v allocations", allocations)

		for inx := uint32(0); inx < allocations; inx++ {
			prov.slots <- 1
		}
	}()

	go func() {
		for range allocRecv {
			log.Printf("Allocate: %v allocate", simId)

			err := stream.Send(&pb.AllocateSlot{})
			if err != nil {
				log.Printf("Allocate: send error %v", err)
				break
			}

			mu.Lock()
			allocations++
			mu.Unlock()
		}
	}()

	for {
		_, err := stream.Recv()
		if err != nil {
			log.Printf("Allocate: recv error %v", err)
			break
		}

		log.Printf("Allocate: %v free", simId)

		prov.slots <- 1

		mu.Lock()
		allocations--
		mu.Unlock()
	}

	return
}
