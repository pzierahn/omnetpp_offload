package main

import (
	"flag"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/provider"
)

var clean bool

func init() {
	flag.BoolVar(&clean, "clean", false, "clean cache")
}

func main() {

	gconfig.ParseFlags()

	if clean {
		provider.Clean()
		return
	}

	//ch := make(chan os.Signal)
	//signal.Notify(ch, os.Interrupt)
	//
	//go func() {
	//	<-ch
	//	provider.Clean()
	//	os.Exit(0)
	//}()

	provider.Start()
}
