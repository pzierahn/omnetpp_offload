package main

import (
	"github.com/pzierahn/project.go.omnetpp/broker"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"os"
	"strconv"
)

func main() {

	config := gconfig.ParseFlagsBroker()

	port := os.Getenv("PORT")
	if port != "" {
		config.BrokerPort, _ = strconv.Atoi(port)
	}

	if err := broker.Start(config); err != nil {
		panic(err)
	}
}
