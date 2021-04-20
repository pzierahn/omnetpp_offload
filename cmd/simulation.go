package main

import (
	"com.github.patrickz98.omnet/simulation"
	"flag"
	"fmt"
)

var simulate string
var simulationName string

func init() {
	flag.StringVar(&simulate, "run", "", "path to OMNeT++ simulation")
	flag.StringVar(&simulationName, "name", "", "name of OMNeT++ simulation")
}

func main() {

	flag.Parse()

	if simulate == "" {
		fmt.Println("missing parameter: run")
		return
	}

	config := simulation.New(simulate, simulationName)

	simulation.Run(config)
}
