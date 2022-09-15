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

var (
	name       = flag.String("scenario", "", "set scenario name")
	broker     = flag.String("broker", "", "set broker addr")
	repeat     = flag.Int("repeat", 5, "repeat trail")
	connect    = flag.String("connect", "", "connect p2p,local,relay")
	simulation = flag.String("simulation", "", "path to simulation")
)

func readSimulationConfig() (runConfig *consumer.Config) {

	simConfPath := filepath.Join(*simulation, "opp-offload-config.json")
	byt, err := os.ReadFile(simConfPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(byt, &runConfig)
	if err != nil {
		log.Fatalln(err)
	}

	runConfig.Path = *simulation
	runConfig.Connect = stargrpc.NameToConnection(*connect)

	return
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	simConfig := readSimulationConfig()

	for trail := 0; trail < *repeat; trail++ {
		log.Printf("Running evaluation %s ==> %d", *name, trail)

		simConfig.Scenario = *name
		simConfig.Trail = fmt.Sprint(trail)

		ctx := context.Background()
		consumer.OffloadSimulation(ctx, gconfig.Broker{
			Address:      *broker,
			BrokerPort:   gconfig.DefaultBrokerPort,
			StargatePort: stargate.DefaultPort,
		}, simConfig)
	}
}
