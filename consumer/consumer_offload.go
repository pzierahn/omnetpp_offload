package consumer

import (
	"context"
	"github.com/pzierahn/omnetpp_offload/eval"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	"github.com/pzierahn/omnetpp_offload/simple"
	"github.com/pzierahn/omnetpp_offload/stargate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// OffloadSimulation starts the simulation offloading to providers.
func OffloadSimulation(ctx context.Context, bconfig gconfig.Broker, config *Config) {

	start := time.Now()

	stargate.SetConfig(stargate.Config{
		Addr: bconfig.Address,
		Port: bconfig.StargatePort,
	})

	if config.Tag == "" {
		config.Tag = filepath.Base(config.Path)
	}

	id := simple.NamedId(config.Tag, 8)
	log.Printf("OffloadSimulation: simulationId %s", id)
	log.Printf("OffloadSimulation: connecting to broker (%v)", bconfig.BrokerDialAddr())

	conn, err := grpc.Dial(
		bconfig.BrokerDialAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		log.Printf("closing connection to broker")
		_ = conn.Close()
	}()

	scenario, trail := os.Getenv("SCENARIO"), os.Getenv("TRAIL")
	if scenario != "" && trail != "" {
		//
		// Enable evaluation logging
		//

		eval.Init(conn, "")

		eval.Start(ctx, scenario, trail, id)
		eval.LogDevice("", 0)

		defer eval.Finish()
	}

	log.Printf("OffloadSimulation: zipping %s", config.Path)

	finish := eval.Log(eval.Event{
		Activity: eval.ActivityCompress,
		Filename: "source.tgz",
	})

	buf, err := simple.TarGz(config.Path, id, config.Ignore...)
	finish(err, 0)

	if err != nil {
		log.Fatalln(err)
	}

	onInit := make(chan uint32)
	defer close(onInit)

	sim := &simulation{
		id:       id,
		ctx:      ctx,
		config:   config,
		queue:    newQueue(),
		source:   buf.Bytes(),
		archLock: make(map[string]*sync.Mutex),
		binaries: make(map[string][]byte),
		onInit:   onInit,
	}

	go sim.startConnector(conn)

	sim.finished.Add(int(<-onInit))
	sim.finished.Wait()
	sim.queue.close()

	log.Printf("OffloadSimulation: simulation %s finished in %v!", id, time.Now().Sub(start))

	return
}
