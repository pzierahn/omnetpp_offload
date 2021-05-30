package main

import (
	"flag"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	"github.com/patrickz98/project.go.omnetpp/provider"
	"os"
	"os/signal"
)

var config gconfig.Config
var clean bool

func init() {
	flag.BoolVar(&clean, "clean", false, "clean cache dir")
	config = gconfig.SourceAndParse()
}

func main() {

	if clean {
		provider.Clean()
		return
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		provider.Clean()
		os.Exit(0)
	}()

	provider.Start(config)

	//conn, err := provider.Init(config)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("############# ccc")
	//
	//if err = conn.StartLink(context.Background()); err != nil {
	//	panic(err)
	//}
	//
	//if err = conn.Close(); err != nil {
	//	panic(err)
	//}
}
