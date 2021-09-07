package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	repeat = 3
)

var writer *csv.Writer
var simulation string

func base() {
	for inx := 0; inx < repeat; inx++ {
		cmd := exec.Command("opp_runall", "-j", "8", "./tictoc", "-c", "TicToc18")
		cmd.Dir = simulation

		log.Printf("Run 0 --> %d", inx)
		start := time.Now()
		if err := cmd.Run(); err != nil {
			panic(err)
		}
		end := time.Now()

		duration := end.Sub(start)
		_ = writer.Write([]string{
			"0",
			fmt.Sprint(inx),
			duration.String(),
		})

		writer.Flush()
	}
}

func scenario(scenario string) {
	for inx := 0; inx < repeat; inx++ {
		cmd := exec.Command("opp_edge_run", "-broker", "85.214.35.83")
		cmd.Dir = simulation
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "SCENARIOID="+scenario)
		cmd.Env = append(cmd.Env, fmt.Sprintf("TRAILID=%d", inx))

		log.Printf("Run %s --> %d", scenario, inx)
		start := time.Now()
		if err := cmd.Run(); err != nil {
			panic(err)
		}
		end := time.Now()

		duration := end.Sub(start)
		_ = writer.Write([]string{
			scenario,
			fmt.Sprint(inx),
			duration.String(),
		})

		writer.Flush()

		_ = os.RemoveAll(filepath.Join(simulation, "opp-edge-results"))
	}
}

func main() {

	var scenarioId string
	flag.StringVar(&scenarioId, "s", "", "scenario")
	flag.Parse()

	if runtime.GOOS == "darwin" {
		simulation = "/Users/patrick/Desktop/tictoc"
	} else {
		simulation = "/home/fioo/patrick/tictoc"
	}

	var filename string

	if scenarioId == "" {
		dir := "system-overhead-scenarios"
		_ = os.MkdirAll(dir, 0755)
		filename = filepath.Join(dir, "system-overhead-local.csv")
	} else {
		filename = fmt.Sprintf("system-overhead-%s.csv", scenarioId)
	}

	_ = os.Remove(filename)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	writer = csv.NewWriter(file)
	defer writer.Flush()

	_ = writer.Write([]string{"scenario", "run", "duration"})

	if scenarioId == "" {
		//
		// Local
		//

		log.Println("Local")

		base()

		ctx, cnl := context.WithCancel(context.Background())
		worker := exec.CommandContext(ctx, "opp_edge_worker")
		if err := worker.Start(); err != nil {
			panic(err)
		}

		scenario("1")

		cnl()
	} else {
		log.Println("Record scenario: " + scenarioId)
		scenario(scenarioId)
	}

	////
	//// Local with opp_edge and docker
	////
	//
	//docker := exec.Command(
	//	"docker", "run", "--rm", "-d",
	//	"pzierahn/omnetpp_edge", "opp_edge_worker", "-broker", "85.214.35.83")
	//
	//id, err := docker.CombinedOutput()
	//if err != nil {
	//	panic(err)
	//}
	//
	//scenario("1")
	//
	//log.Printf("############### docker kill %s", string(id))
	//
	////kill := exec.Command("docker", "kill", string(id))
	////kill.Env = os.Environ()
	////if byt, err := kill.CombinedOutput(); err != nil {
	////	fmt.Println(string(byt))
	////	panic(err)
	////}
}
