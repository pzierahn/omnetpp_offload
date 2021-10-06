package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/scenario"
	"google.golang.org/grpc/benchmark/flags"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const (
	broker = "85.214.35.83"
)

var (
	defaultJobNums = []int{1, 2, 4, 6, 8}
	start          int
	repeat         = 5
	docker         bool
	local          bool
	startWorker    bool
	jobNums        []int
	connects       []string
	scenarioName   string
)

var cancel context.CancelFunc

func init() {
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
}

func csvFile(scenario string) (file *os.File, writer *csv.Writer) {
	dir := filepath.Join("evaluation", "meta", scenario)
	_ = os.MkdirAll(dir, 0755)

	filename := filepath.Join(dir, "durations.csv")

	_ = os.Remove(filename)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("unable to create logfile: %s", err)
	}

	writer = csv.NewWriter(file)
	return
}

func runEvaluation(runner scenario.Runner, connect string, jobs int) error {

	name := scenarioName
	if name == "" {
		name = fmt.Sprintf("c%sj%d", string(connect[0]), jobs)
		if docker {
			name += "d"
		}

		//if local {
		//	name += "l"
		//}
	}

	file, writer := csvFile(name)
	defer func() {
		writer.Flush()
		_ = file.Close()
	}()

	_ = writer.Write([]string{"scenarioId", "trailId", "duration"})
	defer writer.Flush()

	for inx := start; inx < repeat; inx++ {
		duration, err := runner.RunScenario(name, connect, inx)

		if err != nil {
			return err
		}

		record := []string{
			name, fmt.Sprint(inx), duration.String(),
		}

		log.Printf("record: %v", record)
		_ = writer.Write(record)
		writer.Flush()

		// Wait some time to clear buffers.
		time.Sleep(time.Second * 3)
	}

	return nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	nums := flags.IntSlice("j", defaultJobNums, "set parallel job numbers")
	conns := flags.StringSlice("c", []string{"p2p", "relay"}, "set connection")
	flag.IntVar(&start, "ts", 0, "trail start")
	flag.IntVar(&repeat, "tr", repeat, "repeat")
	flag.BoolVar(&docker, "d", false, "use docker")
	flag.BoolVar(&local, "l", false, "start from local")
	flag.BoolVar(&startWorker, "w", false, "start worker")
	flag.StringVar(&scenarioName, "n", "", "scenario name")
	flag.Parse()

	jobNums = *nums
	connects = *conns

	log.Printf("###################################################")
	log.Printf("jobNums: %v", jobNums)
	log.Printf("connects: %v", connects)
	log.Printf("start: %v", start)
	log.Printf("repeat: %v", repeat)
	log.Printf("docker: %v", docker)
	log.Printf("local: %v", local)
	log.Printf("startWorker: %v", startWorker)
	log.Printf("scenarioName: %v", scenarioName)
	log.Printf("###################################################")

	worker := scenario.NewWorker(broker)

	var runner scenario.Runner

	if local {
		home, _ := os.UserHomeDir()
		runner = scenario.NewScenario(scenario.Simulation{
			Broker:     broker,
			OppEdge:    filepath.Join(home, "/github/project.go.omnetpp"),
			Simulation: filepath.Join(home, "/github/project.go.omnetpp/evaluation/tictoc"),
		})
	} else {
		runner = scenario.NewScenarioRemote(scenario.Simulation{
			Broker:     broker,
			OppEdge:    "~/patrick/project.go.omnetpp",
			Simulation: "~/patrick/project.go.omnetpp/evaluation/tictoc",
		})
	}

	if !startWorker {
		for _, connect := range connects {
			if err := runEvaluation(runner, connect, 0); err != nil {
				cancel()
				log.Fatalln(err)
			}
		}

		return
	}

	//defer cancel()

	for _, jobs := range jobNums {
		if docker {
			cancel = worker.StartDocker(jobs)
		} else {
			cancel = worker.StartNative(jobs)
		}

		time.Sleep(time.Second * 3)

		for _, connect := range connects {
			if err := runEvaluation(runner, connect, jobs); err != nil {
				cancel()
				log.Fatalln(err)
			}
		}

		cancel()
	}
}
