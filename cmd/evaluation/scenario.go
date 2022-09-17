package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/consumer"
	"github.com/pzierahn/omnetpp_offload/eval"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	"github.com/pzierahn/omnetpp_offload/stargate"
	"github.com/pzierahn/omnetpp_offload/stargrpc"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

type Worker struct {
	Name   string `json:"name"`
	Docker bool   `json:"docker"`
	Jobs   int    `json:"jobs"`
}

var (
	name       = flag.String("scenario", "", "set scenario name")
	broker     = flag.String("broker", "", "set broker addr")
	repeat     = flag.Int("repeat", 5, "repeat trail")
	connect    = flag.String("connect", "", "connect p2p,local,relay")
	simulation = flag.String("simulation", "", "path to simulation")
	worker     = flag.String("worker", "", "path to worker JSON")
)

func readSimulationConfig() (runConfig *consumer.Config) {

	simConfPath := filepath.Join(*simulation, "opp-offload-config.json")
	byt, err := os.ReadFile(simConfPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(byt, &runConfig)
	if err != nil {
		log.Fatalln(err)
	}

	runConfig.Path = *simulation
	runConfig.Connect = stargrpc.NameToConnection(*connect)

	return
}

func readWorkers() (workers []Worker) {

	if *worker == "" {
		return
	}

	byt, err := os.ReadFile(*worker)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(byt, &workers)
	if err != nil {
		log.Fatalln(err)
	}

	return
}

func killOnExit(stopFuncs []context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		select {
		case rec := <-sig:
			log.Printf("received Interrupt %v\n", rec)

			for _, cnl := range stopFuncs {
				cnl()
			}

			os.Exit(1)
		}
	}()
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	workers := readWorkers()
	var stopFuncs []context.CancelFunc

	for _, start := range workers {

		workerConf := eval.WorkerConfig{
			Name:   start.Name,
			Broker: *broker,
			Jobs:   start.Jobs,
		}

		var cnl context.CancelFunc
		if start.Docker {
			cnl = eval.StartDockerWorker(workerConf)
		} else {
			cnl = eval.StartWorker(workerConf)
		}

		stopFuncs = append(stopFuncs, cnl)
	}

	killOnExit(stopFuncs)

	simConfig := readSimulationConfig()

	if len(workers) > 0 {
		log.Printf("Sleep 8 seconds to ensure all workers have started")
		time.Sleep(time.Second * 8)
	}

	for trail := 0; trail < *repeat; trail++ {
		log.Printf("Running evaluation %s ==> %d", *name, trail)

		simConfig.Scenario = *name
		simConfig.Trail = fmt.Sprint(trail)

		ctx := context.Background()
		consumer.OffloadSimulation(ctx, gconfig.Broker{
			Address:      *broker,
			BrokerPort:   gconfig.DefaultBrokerPort,
			StargatePort: stargate.DefaultPort,
		}, simConfig)
	}
}
