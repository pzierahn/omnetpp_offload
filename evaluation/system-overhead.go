package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	repeat = 3
)

func main() {

	filename := "system-overhead.csv"
	_ = os.Remove(filename)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	_ = writer.Write([]string{"run", "duration", "scenario"})

	//
	// Local
	//

	for inx := 0; inx < repeat; inx++ {
		cmd := exec.Command("opp_runall", "-j", "8", "./tictoc", "-c", "TicToc18")
		cmd.Dir = "/Users/patrick/Desktop/tictoc"
		cmd.Env = []string{
			"ScenarioId", "1",
			"TrailId", fmt.Sprint(inx),
		}

		log.Printf("Run %d", inx)
		start := time.Now()
		if err := cmd.Run(); err != nil {
			panic(err)
		}
		end := time.Now()

		duration := end.Sub(start)
		_ = writer.Write([]string{
			fmt.Sprint(inx),
			duration.String(),
			"opp_runall",
		})
	}

	//
	// Local with opp_edge
	//

	ctx, cnl := context.WithCancel(context.Background())
	worker := exec.CommandContext(ctx, "opp_edge_worker")
	if err := worker.Start(); err != nil {
		panic(err)
	}

	for inx := 0; inx < repeat; inx++ {
		cmd := exec.Command("opp_edge_run")
		cmd.Dir = "/Users/patrick/Desktop/tictoc"
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "SCENARIOID=2")
		cmd.Env = append(cmd.Env, fmt.Sprintf("TRAILID=%d", inx))

		log.Printf("Run %d %s", inx, cmd.Path)
		start := time.Now()
		if err := cmd.Run(); err != nil {
			panic(err)
		}
		end := time.Now()

		duration := end.Sub(start)
		_ = writer.Write([]string{
			fmt.Sprint(inx),
			duration.String(),
			"opp_edge_run",
		})
	}

	cnl()

	//
	// Local with opp_edge and docker
	//

	docker := exec.Command(
		"docker", "run", "--rm", "-d",
		"pzierahn/omnetpp_edge", "opp_edge_worker", "-broker", "85.214.35.83")

	id, err := docker.CombinedOutput()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 3)

	for inx := 0; inx < repeat; inx++ {
		cmd := exec.Command("opp_edge_run")
		cmd.Dir = "/Users/patrick/Desktop/tictoc"
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "SCENARIOID=3")
		cmd.Env = append(cmd.Env, fmt.Sprintf("TRAILID=%d", inx))

		log.Printf("Run %d %s", inx, cmd.Path)
		start := time.Now()
		if err := cmd.Run(); err != nil {
			panic(err)
		}
		//if err := cmd.Run(); err != nil {
		//	panic(err)
		//}
		end := time.Now()

		duration := end.Sub(start)
		_ = writer.Write([]string{
			fmt.Sprint(inx),
			duration.String(),
			"opp_edge_run docker",
		})
	}

	log.Printf("############### docker kill %s", string(id))

	//kill := exec.Command("docker", "kill", string(id))
	//kill.Env = os.Environ()
	//if byt, err := kill.CombinedOutput(); err != nil {
	//	fmt.Println(string(byt))
	//	panic(err)
	//}
}
