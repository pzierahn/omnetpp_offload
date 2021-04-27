package main

import (
	"context"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	"github.com/patrickz98/project.go.omnetpp/worker"
)

var config gconfig.Config

func init() {
	config = gconfig.SourceAndParse()
}

func main() {

	conn, err := worker.Init(config)
	if err != nil {
		panic(err)
	}

	if err = conn.StartLink(context.Background()); err != nil {
		panic(err)
	}
}
