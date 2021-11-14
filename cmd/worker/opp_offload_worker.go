package main

import (
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

	// TODO: Cloud compatibility.
	//envPort := os.Getenv("PORT")
	//if envPort != "" {
	//	config.Port, _ = strconv.Atoi(envPort)
	//}

	provider.Start(config)
}
