package main

import (
	"com.github.patrickz98.omnet/simulation"
	"flag"
)

var status bool
var simulate string

func init() {
	flag.BoolVar(&status, "status", false, "")
	flag.StringVar(&simulate, "simulate", "", "")
}

func main() {

	flag.Parse()

	if simulate != "" {
		simulation.Run(simulate)
	}

	if status {
		simulation.Status()
	}
}
