package main

import (
	"flag"
	"github.com/pzierahn/project.go.omnetpp/broker"
	"github.com/pzierahn/project.go.omnetpp/defines"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"os"
	"os/signal"
)

var (
	config broker.Config
	clean  bool
)

func init() {
	flag.BoolVar(&clean, "clean", false, "clean broker")
	flag.IntVar(&config.BrokerPort, "port", defines.DefaultPort, "set broker port")
}

func main() {

	flag.Parse()

	if clean {
		storage.Clean()
		return
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		storage.Clean()
		os.Exit(0)
	}()

	if err := broker.Start(config); err != nil {
		panic(err)
	}
}
