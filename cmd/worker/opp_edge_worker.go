package main

import (
	"flag"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/provider"
	"github.com/pzierahn/project.go.omnetpp/storage"
)

var clean bool

func init() {
	flag.BoolVar(&clean, "clean", false, "clean cache")
}

func main() {

	gconfig.ParseFlags()

	if clean {
		provider.Clean()
		storage.Clean()
		return
	}

	provider.Start()
}
