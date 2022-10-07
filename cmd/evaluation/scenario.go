package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pzierahn/omnetpp_offload/consumer"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	"github.com/pzierahn/omnetpp_offload/provider"
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

type workerConfig struct {
	name   string
	broker string
	jobs   int
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

func startDockerWorker(worker workerConfig) (cancel context.CancelFunc) {
	log.Printf("Starting docker worker: worker=%+v", worker)

	var err error
	dockerClI, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("unable to connect to docker: %s", err)
	}

	ctx := context.Background()
	resp, err := dockerClI.ContainerCreate(ctx, &container.Config{
		Image: "pzierahn/omnetpp_offload",
		Cmd: []string{
			"opp_offload_worker",
			"-broker", worker.broker,
			"-name", worker.name,
			"-jobs", fmt.Sprint(worker.jobs),
		},
	}, &container.HostConfig{
		SecurityOpt: []string{
			"seccomp:unconfined",
		},
	}, nil, nil, worker.name)
	if err != nil {
		log.Panic(err)
	}

	if err := dockerClI.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return func() {
		log.Printf("Stopping container %s", resp.ID)
		if err := dockerClI.ContainerStop(ctx, resp.ID, nil); err != nil {
			panic(err)
		}

		log.Printf("Removing container %s", resp.ID)
		if err := dockerClI.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
			panic(err)
		}
	}

	return
}

func startWorker(worker workerConfig) (cancel context.CancelFunc) {
	log.Printf("Starting worker: worker=%+v", worker)

	ctx, cnl := context.WithCancel(context.Background())

	go func() {
		provider.Start(ctx, gconfig.Config{
			Provider: gconfig.Provider{
				Name: worker.name,
				Jobs: worker.jobs,
			},
			Broker: gconfig.Broker{
				Address:      worker.broker,
				BrokerPort:   gconfig.DefaultBrokerPort,
				StargatePort: stargate.DefaultPort,
			},
		})
	}()

	return func() {
		log.Printf("Stopping worker %s", worker.name)
		cnl()
	}
}
func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	workers := readWorkers()
	var stopFuncs []context.CancelFunc

	for _, start := range workers {

		workerConf := workerConfig{
			name:   start.Name,
			broker: *broker,
			jobs:   start.Jobs,
		}

		var cnl context.CancelFunc
		if start.Docker {
			cnl = startDockerWorker(workerConf)
		} else {
			cnl = startWorker(workerConf)
		}

		stopFuncs = append(stopFuncs, cnl)
	}

	killOnExit(stopFuncs)

	simConfig := readSimulationConfig()

	if len(workers) > 0 {
		log.Printf("Sleep 4 seconds to ensure all workers have started")
		time.Sleep(time.Second * 4)
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
