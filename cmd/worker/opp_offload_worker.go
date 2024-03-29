package main

import (
	"context"
	"flag"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	"github.com/pzierahn/omnetpp_offload/provider"
	"github.com/pzierahn/omnetpp_offload/storage"
)

var clean bool

func init() {
	flag.BoolVar(&clean, "clean", false, "clean all cache files")
}

func main() {

	config := gconfig.ParseFlags()

	if clean {
		provider.Clean()
		storage.Clean()
		return
	}

	provider.Start(context.Background(), config)
}
