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
	"os/exec"
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
	cancel    context.CancelFunc
)

var (
	defaultJobNums  = []int{1, 2, 4, 6, 8}
	defaultConnects = []string{"p2p", "relay"}
	start           int
	repeat          = 5
	docker          bool
	jobNums         []int
	connects        []string
)

func initDockerSSH() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		select {
		case rec := <-sig:
			log.Printf("received Interrupt %v\n", rec)
			cancel()
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

	start := time.Now()
	if err = session.Run(strings.Join(cmd, "; ")); err != nil {
		err = fmt.Errorf("unable to run scenario: %s", err)
		return
	}

	duration = time.Now().Sub(start)

	return
}

func startDocker(jobs int) (cancel context.CancelFunc) {
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

	return func() {
		stop(resp.ID)
	}
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

func startNative(jobs int) (cancel context.CancelFunc) {

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "go", "run", "cmd/worker/opp_edge_worker.go",
		"-broker", broker,
		"-jobs", fmt.Sprint(jobs))
	if err := cmd.Start(); err != nil {
		log.Fatalf("unable to start worker: %s", err)
	}

	return
}

func overheadFile(scenario string) (file *os.File, err error) {

	dir := "evaluation/meta"
	_ = os.MkdirAll(dir, 0755)

	filename := filepath.Join(dir, fmt.Sprintf("overhead-%s.csv", scenario))

	_ = os.Remove(filename)

	return os.Create(filename)
}

func runEvaluation(connect string, jobs int) error {

	var scenario string
	if docker {
		scenario = fmt.Sprintf("%sj%dd", string(connect[0]), jobs)
	} else {
		scenario = fmt.Sprintf("%sj%d", string(connect[0]), jobs)
	}

	file, err := overheadFile(scenario)
	if err != nil {
		return fmt.Errorf("unable to create logfile: %s", err)
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	_ = writer.Write([]string{"scenarioId", "trailId", "duration"})

	for inx := start; inx < repeat; inx++ {
		duration, err := runScenario(scenario, connect, inx)

		if err != nil {
			return err
		}

		record := []string{
			scenario, fmt.Sprint(inx), duration.String(),
		}

		log.Printf("record: %v", record)
		_ = writer.Write(record)
		writer.Flush()

		// Wait some time to clear buffers.
		time.Sleep(time.Second * 2)
	}

	return nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	nums := flags.IntSlice("j", defaultJobNums, "set parallel job numbers")
	conns := flags.StringSlice("c", defaultConnects, "set connection")
	flag.IntVar(&start, "ts", 0, "trail start")
	flag.IntVar(&repeat, "tr", repeat, "repeat")
	flag.BoolVar(&docker, "d", false, "use docker")
	flag.Parse()

	jobNums = *nums
	connects = *conns

	log.Printf("jobNums: %v", jobNums)
	log.Printf("connects: %v", connects)
	log.Printf("start: %v", start)
	log.Printf("repeat: %v", repeat)
	log.Printf("docker: %v", docker)

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

	for _, jobs := range jobNums {
		if docker {
			cancel = startDocker(jobs)
		} else {
			cancel = startNative(jobs)
		}

		time.Sleep(time.Second * 3)

		for _, connect := range connects {
			if err := runEvaluation(connect, jobs); err != nil {
				cancel()
				log.Fatalln(err)
			}
		}

		cancel()
	}
}
