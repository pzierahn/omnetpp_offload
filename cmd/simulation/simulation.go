package main

import (
	"com.github.patrickz98.omnet/simulation"
	"flag"
	"fmt"
)

var debug bool
var simulate string
var name string

func init() {
	flag.BoolVar(&debug, "debug", false, "send debug request")
	flag.StringVar(&simulate, "run", "", "path to OMNeT++ simulation")
	flag.StringVar(&name, "name", "", "name of OMNeT++ simulation")
}

func main() {

	flag.Parse()

	if debug {
		simulation.DebugRequest()
		return
	}

	if simulate == "" {
		fmt.Println("missing parameter: run")
		return
	}

	config := simulation.New(simulate, name)

	simulation.Run(config)
}
