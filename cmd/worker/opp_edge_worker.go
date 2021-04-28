package main

import (
	"context"
	"flag"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	"github.com/patrickz98/project.go.omnetpp/worker"
)

var config gconfig.Config
var clean bool

func init() {
	flag.BoolVar(&clean, "clean", false, "clean cache dir")
	config = gconfig.SourceAndParse()
}

func main() {

	if clean {
		worker.Clean()
	}

	conn, err := worker.Init(config)
	if err != nil {
		panic(err)
	}

	if err = conn.StartLink(context.Background()); err != nil {
		panic(err)
	}
}
