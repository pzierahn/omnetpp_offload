package main

import (
	"context"
	"flag"
	"github.com/pzierahn/omnetpp_offload/stargate"
	"log"
	"os"
	"strconv"
)

var port int

func init() {
	flag.IntVar(&port, "port", stargate.DefaultPort, "set stargate server")
	flag.Parse()
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := stargate.Config{
		Port: port,
	}

	envPort := os.Getenv("PORT")
	if envPort != "" {
		config.Port, _ = strconv.Atoi(envPort)
	}

	// Set stun server address
	stargate.SetConfig(config)

	err := stargate.Server(context.Background(), true)
	if err != nil {
		log.Fatalln(err)
	}
}
