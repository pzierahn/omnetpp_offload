package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	projectDir = "~/patrick/project.go.omnetpp"
)

var (
	cli       *client.Client
	sshClient *ssh.Client
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	connectSSH()
}

func connectSSH() {
	var err error

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile(filepath.Join(home, "/.ssh/id_rsa"))
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: "fioo",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
			//ssh.Password("PASSWORD"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err = ssh.Dial("tcp", "dc1.fioo.one:4777", config)
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}
}

func updateRepo() {

	log.Printf("Updating repository")
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatalf("unable to create a new session: %s", err)
	}
	defer func() {
		_ = session.Close()
	}()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run(fmt.Sprintf("cd %s; git pull", projectDir)); err != nil {
		log.Fatalf("unable to update project: %s", err)
	}
}

func runScenario(name, connect string, repeat int) {

	log.Printf("\n\n##################################################")
	log.Printf("Staring scenario %s %s, %d", name, connect, repeat)
	log.Printf("##################################################\n\n")

	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatalf("unable to create a new session: %s", err)
	}
	defer func() {
		_ = session.Close()
	}()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	//cmd := "cd " + projectDir + ";" +
	//	"go run cmd/scenario/scenario.go -s "

	cmd := []string{
		fmt.Sprintf("cd %s", projectDir),
		fmt.Sprintf("go run cmd/scenario/scenario.go -s %s -c %s -r %d", name, connect, repeat),
	}

	if err := session.Run(strings.Join(cmd, "; ")); err != nil {
		log.Fatalf("unable to run scenario: %s", err)
	}
}

func startDocker(jobs int) (id string) {
	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "pzierahn/omnetpp_edge",
		Cmd: []string{
			"opp_edge_worker",
			"-broker", "85.214.35.83",
			"-name", "docker",
			"-jobs", fmt.Sprint(jobs),
		},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}

func stop(id string) {
	log.Printf("Stopping container %s", id)
	if err := cli.ContainerStop(context.Background(), id, nil); err != nil {
		panic(err)
	}

	log.Printf("Removing container %s", id)
	if err := cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
}

func main() {

	defer func() {
		_ = cli.Close()
		_ = sshClient.Close()
	}()

	updateRepo()

	jobNums := []int{1, 2, 4, 6, 8}
	connects := []string{
		"p2p", "relay",
	}

	for _, connect := range connects {
		for _, jobs := range jobNums {
			id := startDocker(jobs)
			// Sleep to ensure that the docker is started and connected.
			time.Sleep(time.Second * 3)

			name := fmt.Sprintf("%sj%d", string(connect[0]), jobs)
			runScenario(name, connect, 3)

			stop(id)
		}
	}
}
