package main

import (
	"com.github.patrickz98.omnet/omnetpp"
	"flag"
	"fmt"
)

var path string

func init() {
	flag.StringVar(&path, "path", "", "simulation path")
}

func main() {
	flag.Parse()

	if path == "" {
		fmt.Println("path to source is missing!")
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
