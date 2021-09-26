package consumer

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/eval"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"google.golang.org/grpc"
	"log"
	"path/filepath"
	"sync"
)

// OffloadSimulation starts the simulation offloading to providers.
func OffloadSimulation(ctx context.Context, config *Config) {

	if config.Tag == "" {
		config.Tag = filepath.Base(config.Path)
	}

	id := simple.NamedId(config.Tag, 8)
	log.Printf("OffloadSimulation: simulationId %s", id)
	log.Printf("OffloadSimulation: connecting to broker (%v)", gconfig.BrokerDialAddr())

	conn, err := grpc.Dial(
		gconfig.BrokerDialAddr(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	eval.Init(conn)
	eval.SetScenario(id)
	eval.DeviceId = "consumer"

	log.Printf("OffloadSimulation: zipping %s", config.Path)

	done := eval.LogAction(eval.ActionCompress, "source")
	buf, err := simple.TarGz(config.Path, id, config.Ignore...)
	_ = done(err)
	if err != nil {
		log.Fatalln(err)
	}

	sim := &simulation{
		id:       id,
		ctx:      ctx,
		config:   config,
		queue:    newQueue(),
		source:   buf.Bytes(),
		archLock: make(map[string]*sync.Mutex),
		binaries: make(map[string][]byte),
	}

	onInit := make(chan uint32)
	defer close(onInit)

	go sim.startConnector(conn, onInit)

	sim.finished.Add(int(<-onInit))
	sim.finished.Wait()

	//time.Sleep(time.Second*3)
	//log.Printf("########### Add somemore stuff.")
	//
	//sim.finished.Add(4)
	//
	//sim.queue.add(&pb.SimulationRun{
	//	SimulationId: id,
	//	Config:       "TicToc18",
	//	RunNum:       "1",
	//}, &pb.SimulationRun{
	//	SimulationId: id,
	//	Config:       "TicToc18",
	//	RunNum:       "2",
	//}, &pb.SimulationRun{
	//	SimulationId: id,
	//	Config:       "TicToc18",
	//	RunNum:       "3",
	//}, &pb.SimulationRun{
	//	SimulationId: id,
	//	Config:       "TicToc18",
	//	RunNum:       "4",
	//})
	//
	//sim.finished.Wait()

	sim.queue.killLingering()

	log.Printf("OffloadSimulation: simulation %s finished!", id)

	return
}
