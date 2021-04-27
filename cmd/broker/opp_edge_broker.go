package main

import (
	"flag"
	"github.com/patrickz98/project.go.omnetpp/broker"
	"github.com/patrickz98/project.go.omnetpp/defines"
)

var (
	config broker.Config
)

func init() {
	flag.IntVar(&config.Port, "port", defines.DefaultPort, "set broker port")
}

func main() {

	flag.Parse()

	if err := broker.Start(config); err != nil {
		panic(err)
	}
}
