package main

import (
	"github.com/patrickz98/project.go.omnetpp/stargate"
	"os"
)

func main() {
	cmd := os.Args[1]
	switch cmd {
	case "c":
		stargate.Client()
	case "s":
		stargate.Server()
	}
}
