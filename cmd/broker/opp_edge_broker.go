package main

import (
	"github.com/pzierahn/project.go.omnetpp/broker"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
)

func main() {

	gconfig.ParseFlags()

	if err := broker.Start(); err != nil {
		panic(err)
	}
}
