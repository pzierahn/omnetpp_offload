package main

import (
	"github.com/patrickz98/project.go.omnetpp/broker"
)

func main() {

	if err := broker.Start(); err != nil {
		panic(err)
	}
}
