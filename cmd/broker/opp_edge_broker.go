package main

import (
	"flag"
	"github.com/patrickz98/project.go.omnetpp/broker"
	"github.com/patrickz98/project.go.omnetpp/defines"
	"github.com/patrickz98/project.go.omnetpp/storage"
)

var (
	config broker.Config
	clean  bool
)

func init() {
	flag.BoolVar(&clean, "clean", false, "clean broker")
	flag.IntVar(&config.Port, "port", defines.DefaultPort, "set broker port")
}

func main() {

	flag.Parse()

	if clean {
		storage.Clean()
	}

	if err := broker.Start(config); err != nil {
		panic(err)
	}
}
