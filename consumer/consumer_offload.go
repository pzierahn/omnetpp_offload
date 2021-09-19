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
	eval.DeviceId = "simulation"

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

	onInit := make(chan int32)
	defer close(onInit)

	go sim.startConnector(conn, onInit)

	sim.finished.Add(int(<-onInit))
	sim.finished.Wait()

	log.Printf("OffloadSimulation: simulation %s finished!", id)

	return
}
