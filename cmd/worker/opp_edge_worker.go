package main

import (
	"flag"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/provider"
	"github.com/pzierahn/project.go.omnetpp/storage"
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
