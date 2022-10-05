package eval

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	"github.com/pzierahn/omnetpp_offload/provider"
	"github.com/pzierahn/omnetpp_offload/stargate"
	"log"
)

type WorkerConfig struct {
	Name   string
	Broker string
	Jobs   int
}

func StartDockerWorker(worker WorkerConfig) (cancel context.CancelFunc) {
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
			"-broker", worker.Broker,
			"-name", worker.Name,
			"-jobs", fmt.Sprint(worker.Jobs),
		},
	}, nil, nil, nil, worker.Name)
	if err != nil {
		panic(err)
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

func StartWorker(worker WorkerConfig) (cancel context.CancelFunc) {
	log.Printf("Starting worker: worker=%+v", worker)

	ctx, cnl := context.WithCancel(context.Background())

	go func() {
		provider.Start(ctx, gconfig.Config{
			Provider: gconfig.Provider{
				Name: worker.Name,
				Jobs: worker.Jobs,
			},
			Broker: gconfig.Broker{
				Address:      worker.Broker,
				BrokerPort:   gconfig.DefaultBrokerPort,
				StargatePort: stargate.DefaultPort,
			},
		})
	}()

	return func() {
		log.Printf("Stopping worker %s", worker.Name)
		cnl()
	}
}
