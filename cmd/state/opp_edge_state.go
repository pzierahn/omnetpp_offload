package main

import (
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/stateinfo"
)

var config gconfig.Config

func init() {
	config = gconfig.SourceAndParse(gconfig.ParseBroker)
}

func main() {

	stateinfo.Workers(config.Broker)
}
