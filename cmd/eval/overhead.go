package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	repeat = 5
	broker = "85.214.35.83"
	//broker = "localhost"
)

var (
	simulation string
	connection string
)

var writer *csv.Writer

func local() {
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
		_ = os.RemoveAll(filepath.Join(simulation, "opp-edge-results"))
		_ = os.RemoveAll(filepath.Join(simulation, "results"))

		cmd := exec.Command("opp_edge_run", "-broker", broker)
		cmd.Dir = simulation
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "CONNECT="+connection)
		cmd.Env = append(cmd.Env, "SCENARIOID="+scenario)
		cmd.Env = append(cmd.Env, fmt.Sprintf("TRAILID=%d", inx))

		log.Printf("Run scenario: %s trail: %d", scenario, inx)
		start := time.Now()
		simple.RunCmdStdout(cmd)
		//if err := cmd.Run(); err != nil {
		//	panic(err)
		//}
		end := time.Now()

		duration := end.Sub(start)
		_ = writer.Write([]string{
			scenario,
			fmt.Sprint(inx),
			duration.String(),
		})

		writer.Flush()
	}
}

func main() {

	var scenarioId string

	flag.StringVar(&scenarioId, "s", "", "scenario")
	flag.StringVar(&connection, "c", "", "connection: local|p2p|relay")
	flag.Parse()

	if runtime.GOOS == "darwin" {
		simulation = "/Users/patrick/Desktop/tictoc"
	} else {
		simulation = "/home/fioo/patrick/tictoc"
	}

	log.Println("Installing latest opp_edge version...")

	cmd := exec.Command("go", "install", "cmd/consumer/opp_edge_run.go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("out: ", string(out))
		panic(err)
	}

	var filename string

	dir := "evaluation/meta"
	_ = os.MkdirAll(dir, 0755)

	if scenarioId == "" {
		filename = filepath.Join(dir, "overhead-local.csv")
	} else {
		filename = filepath.Join(dir, fmt.Sprintf("overhead-%s.csv", scenarioId))
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
		local()
	} else {
		log.Println("Record scenario: " + scenarioId)
		scenario(scenarioId)
	}
}
