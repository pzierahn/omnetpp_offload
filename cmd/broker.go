package main

import (
	"com.github.patrickz98.omnet/broker"
)

func main() {

	if err := broker.Start(); err != nil {
		panic(err)
	}
}
