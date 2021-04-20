package main

import (
	"com.github.patrickz98.omnet/simulation"
	"flag"
)

var simulate string

func init() {
	flag.StringVar(&simulate, "run", "", "path to OMNeT++ simulation")
}

func main() {

	flag.Parse()

	if simulate != "" {
		simulation.Run(simulate)
	}
}
