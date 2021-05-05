package main

import (
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	"github.com/patrickz98/project.go.omnetpp/stateinfo"
)

var config gconfig.Config

func init() {
	config = gconfig.SourceAndParse(gconfig.ParseBroker)
}

func main() {

	stateinfo.Status(config.Broker, nil)
}
