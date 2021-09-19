package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc/benchmark/flags"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

const (
	projectDir = "~/patrick/project.go.omnetpp"
	simulation = "/home/fioo/patrick/tictoc"
	broker     = "85.214.35.83"
)

var (
	cli       *client.Client
	sshClient *ssh.Client
	dockerId  string
)

var (
	defaultJobNums  = []int{1, 2, 4, 6, 8}
	defaultConnects = []string{"p2p", "relay"}
	repeat          = 5
	jobNums         []int
	connects        []string
)

func initDockerSSH() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		select {
		case rec := <-sig:
			log.Printf("received SIGTERM %v\n", rec)
			stop(dockerId)
			os.Exit(1)
		}
	}()

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

func runScenario(scenario, connect string, trail int) (duration time.Duration, err error) {

	log.Printf("Staring scenario scenario=%s connect=%s trail=%d", scenario, connect, trail)

	session, err := sshClient.NewSession()
	if err != nil {
		err = fmt.Errorf("unable to create a new session: %s", err)
		return
	}
	defer func() {
		_ = session.Close()
	}()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// session.Setenv() doesn't work.
	// see https://vic.demuzere.be/articles/environment-variables-setenv-ssh-golang/
	envVars := []string{
		fmt.Sprintf("SCENARIOID=%s", scenario),
		fmt.Sprintf("CONNECT=%s", connect),
		fmt.Sprintf("TRAILID=%d", trail),
	}
	envPrefix := strings.Join(envVars, " ")

	cmd := []string{
		// Source GO paths.
		fmt.Sprintf("source ~/.profile"),
		// Switch to project dir
		fmt.Sprintf("cd %s", projectDir),
	}

	// Delete deprecated data.
	deleteDirs := []string{
		filepath.Join(simulation, "opp-edge-results"),
		filepath.Join(simulation, "results"),
		filepath.Join(simulation, "out"),
	}

	for _, dir := range deleteDirs {
		cmd = append(cmd, "rm -rf "+dir)
	}

	cmd = append(cmd,
		fmt.Sprintf("%s go run cmd/consumer/opp_edge_run.go -broker %s -path %s -config %s",
			envPrefix, broker, simulation, filepath.Join(simulation, "opp-edge-config.json")),
	)

	log.Printf("bashscript:\n%s\n", strings.Join(cmd, "; \n"))

	start := time.Now()
	if err = session.Run(strings.Join(cmd, "; ")); err != nil {
		err = fmt.Errorf("unable to run scenario: %s", err)
		return
	}

	duration = time.Now().Sub(start)

	return
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
	if id == "" {
		return
	}

	log.Printf("Stopping container %s", id)
	if err := cli.ContainerStop(context.Background(), id, nil); err != nil {
		panic(err)
	}

	log.Printf("Removing container %s", id)
	if err := cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
}

func overheadFile(scenario string) (file *os.File, err error) {

	dir := "evaluation/meta"
	_ = os.MkdirAll(dir, 0755)

	filename := filepath.Join(dir, fmt.Sprintf("overhead-%s.csv", scenario))

	_ = os.Remove(filename)

	return os.Create(filename)
}

func runEvaluation(connect string, jobs int) {
	dockerId = startDocker(jobs)
	defer stop(dockerId)

	// Sleep to ensure that the docker is started and connected.
	time.Sleep(time.Second * 3)

	scenario := fmt.Sprintf("%sj%dd", string(connect[0]), jobs)
	file, err := overheadFile(scenario)
	if err != nil {
		stop(dockerId)
		log.Fatalf("unable to create logfile: %s", err)
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	_ = writer.Write([]string{"scenarioId", "trailId", "duration"})

	for inx := 0; inx < repeat; inx++ {
		duration, err := runScenario(scenario, connect, inx)

		if err != nil {
			stop(dockerId)
			log.Fatalln(err)
		}

		record := []string{
			scenario, fmt.Sprint(inx), duration.String(),
		}

		log.Printf("record: %v", record)
		_ = writer.Write(record)
		writer.Flush()
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	nums := flags.IntSlice("j", defaultJobNums, "set parallel job numbers")
	conns := flags.StringSlice("c", defaultConnects, "set connection")
	flag.IntVar(&repeat, "r", repeat, "repeat")
	flag.Parse()

	jobNums = *nums
	connects = *conns

	log.Printf("jobNums: %v", jobNums)
	log.Printf("connects: %v", connects)
	log.Printf("repeat: %v", repeat)

	initDockerSSH()
	defer func() {
		_ = cli.Close()
		_ = sshClient.Close()
	}()

	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatalf("unable to create a new session: %s", err)
	}
	defer func() {
		_ = session.Close()
	}()

	updateRepo()

	for _, connect := range connects {
		for _, jobs := range jobNums {
			runEvaluation(connect, jobs)
		}
	}
}
