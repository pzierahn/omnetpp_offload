package main

import (
	"github.com/pzierahn/project.go.omnetpp/broker"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"os"
	"os/signal"
)

func main() {

	gconfig.ParseFlags()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		storage.Clean()
		os.Exit(0)
	}()

	if err := broker.Start(); err != nil {
		panic(err)
	}
}
