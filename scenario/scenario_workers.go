package scenario

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"log"
	"os/exec"
)

type Worker struct {
	docker *client.Client
	broker string
}

func NewWorker(broker string) (worker Worker) {
	worker.broker = broker

	var err error
	worker.docker, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("unable to connect to docker: %s", err)
	}

	return
}

func (worker Worker) StartNative(jobs int) (cancel context.CancelFunc) {

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "go", "run", "cmd/worker/opp_edge_worker.go",
		"-broker", worker.broker,
		"-jobs", fmt.Sprint(jobs))
	if err := cmd.Start(); err != nil {
		log.Fatalf("unable to start worker: %s", err)
	}

	return
}

func (worker Worker) StartDocker(jobs int) (cancel context.CancelFunc) {
	ctx := context.Background()
	resp, err := worker.docker.ContainerCreate(ctx, &container.Config{
		Image: "pzierahn/omnetpp_edge",
		Cmd: []string{
			"opp_edge_worker",
			"-broker", worker.broker,
			"-name", "docker",
			"-jobs", fmt.Sprint(jobs),
		},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := worker.docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return func() {
		worker.stop(resp.ID)
	}
}

func (worker Worker) stop(id string) {
	if id == "" {
		return
	}

	log.Printf("Stopping container %s", id)
	if err := worker.docker.ContainerStop(context.Background(), id, nil); err != nil {
		panic(err)
	}

	log.Printf("Removing container %s", id)
	if err := worker.docker.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
}
