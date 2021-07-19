package main

import (
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/provider"
)

func main() {

	gconfig.ParseFlags()

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
