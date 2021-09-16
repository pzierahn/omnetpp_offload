package main

import (
	"flag"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/provider"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"github.com/pzierahn/project.go.omnetpp/storage"
)

var clean bool

func init() {
	flag.BoolVar(&clean, "clean", false, "clean all cache files")
}

func main() {

	gconfig.ParseFlags(gconfig.ParseAll)

	if clean {
		provider.Clean()
		storage.Clean()
		return
	}

	stargate.SetConfig(stargate.Config{
		Addr: gconfig.BrokerAddr(),
		Port: gconfig.StargatePort(),
	})

	provider.Start()
}
