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

func StartDocker(broker string, jobs int) (cancel context.CancelFunc) {
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
			"-broker", broker,
			"-name", "docker",
			"-jobs", fmt.Sprint(jobs),
		},
	}, nil, nil, nil, "")
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

func StartNative(broker string, jobs int) (cancel context.CancelFunc) {

	ctx, cnl := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "go", "run", "cmd/worker/opp_offload_worker.go",
		"-broker", broker,
		"-jobs", fmt.Sprint(jobs))
	if err := cmd.Start(); err != nil {
		log.Fatalf("unable to start worker: %s", err)
	}

	return func() {
		log.Printf("Stopping worker")
		cnl()
	}
}
