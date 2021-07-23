package main

import (
	"github.com/pzierahn/project.go.omnetpp/broker"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/stargate"
)

func main() {

	gconfig.ParseFlags()

	stargate.SetConfig(stargate.Config{
		Addr: gconfig.BrokerAddr(),
		Port: gconfig.StargatePort(),
	})

	if err := broker.Start(); err != nil {
		panic(err)
	}
}
