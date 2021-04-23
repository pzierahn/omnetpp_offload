package main

import (
	"flag"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/omnetpp"
)

var path string

func init() {
	flag.StringVar(&path, "path", "", "simulation path")
}

func main() {
	flag.Parse()

	if path == "" {
		fmt.Println("missing parameter: path")
		return
	}

	opp := omnetpp.New(path)

	err := opp.Clean()
	if err != nil {
		panic(err)
	}

	err = opp.MakeMake()
	if err != nil {
		panic(err)
	}

	err = opp.Compile()
	if err != nil {
		panic(err)
	}
}
