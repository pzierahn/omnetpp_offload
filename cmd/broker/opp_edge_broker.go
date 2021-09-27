package main

import (
	"github.com/pzierahn/project.go.omnetpp/broker"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
)

func main() {

	config := gconfig.ParseFlagsBroker()

	if err := broker.Start(config); err != nil {
		panic(err)
	}
}
