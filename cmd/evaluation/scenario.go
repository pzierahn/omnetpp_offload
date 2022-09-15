package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/consumer"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	"github.com/pzierahn/omnetpp_offload/stargate"
	"github.com/pzierahn/omnetpp_offload/stargrpc"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	ScenarioName   string `json:"scenario-name"`
	Broker         string `json:"broker"`
	Repeat         int    `json:"repeat"`
	Connect        string `json:"connect"`
	SimulationPath string `json:"simulation-path"`
}

var (
	scenarioPath = flag.String("scenario", "", "set scenario JSON path")
)

func readConfig() (config Config) {
	byt, err := os.ReadFile(*scenarioPath)
	if err != nil {
		log.Fatalf("couldn't read scenario JSON: %v", err)
	}

	err = json.Unmarshal(byt, &config)
	if err != nil {
		log.Fatalf("couldn't parse scenario JSON: %v", err)
	}

	return
}

func readSimulationConfig(config Config) (runConfig *consumer.Config) {

	simConfPath := filepath.Join(config.SimulationPath, "opp-offload-config.json")
	byt, err := os.ReadFile(simConfPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(byt, &runConfig)
	if err != nil {
		log.Fatalln(err)
	}

	runConfig.Path = config.SimulationPath
	runConfig.Connect = stargrpc.NameToConnection(config.Connect)

	return
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	config := readConfig()
	simConfig := readSimulationConfig(config)

	for trail := 0; trail < config.Repeat; trail++ {
		log.Printf("Running evaluation %s ==> %d", config.ScenarioName, trail)

		simConfig.Scenario = config.ScenarioName
		simConfig.Trail = fmt.Sprint(trail)

		ctx := context.Background()
		consumer.OffloadSimulation(ctx, gconfig.Broker{
			Address:      config.Broker,
			BrokerPort:   gconfig.DefaultBrokerPort,
			StargatePort: stargate.DefaultPort,
		}, simConfig)
	}
}
