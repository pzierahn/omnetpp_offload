package main

import (
	"com.github.patrickz98.omnet/broker"
	"com.github.patrickz98.omnet/worker"
	"flag"
)

var start string

func init() {
	flag.StringVar(&start, "start", "", "")
}

func main() {

	flag.Parse()

	if start == "server" {
		_ = broker.Start()
	}

	if start == "client" {
		_ = worker.Link()
	}
}
