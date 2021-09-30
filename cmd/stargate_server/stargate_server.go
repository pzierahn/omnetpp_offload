package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"log"
	"net/http"
	"os"
	"strconv"
)

var port int

func init() {
	flag.IntVar(&port, "port", stargate.DefaultPort, "set stargate server")
	flag.Parse()
}

func hello(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			_, _ = fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := stargate.Config{
		Port: port,
	}

	port := os.Getenv("PORT")
	if port != "" {
		config.Port, _ = strconv.Atoi(port)
	}

	// Set stun server address
	stargate.SetConfig(config)

	err := stargate.Server(context.Background(), true)
	if err != nil {
		log.Fatalln(err)
	}
}
