package main

import (
	"flag"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/simulation"
	"strings"
)

var debug bool
var path string
var name string
var configs string

func init() {
	flag.BoolVar(&debug, "debug", false, "send debug request")
	flag.StringVar(&path, "path", "", "path to OMNeT++ simulation")
	flag.StringVar(&name, "name", "", "name of the simulation")
	flag.StringVar(&configs, "configs", "", "simulation config names")
}

func main() {

	flag.Parse()

	if debug {
		simulation.DebugRequest()
		return
	}

	if path == "" {
		fmt.Println("missing parameter: path")
		return
	}

	config := simulation.New(path, name)

	if configs != "" {
		config.Configs = strings.Split(configs, ",")
	}

	simulation.Run(config)
}
