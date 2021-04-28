package main

import (
	"context"
	"flag"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	"github.com/patrickz98/project.go.omnetpp/worker"
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
		worker.Clean()
		return
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		worker.Clean()
		os.Exit(1)
	}()

	conn, err := worker.Init(config)
	if err != nil {
		panic(err)
	}

	if err = conn.StartLink(context.Background()); err != nil {
		panic(err)
	}
}
