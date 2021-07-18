package main

import (
	"flag"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"log"
)

func init() {
	_ = flag.Bool("bool", false, "")
	flag.Parse()
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	addr1, addr2, err := stargate.RelayServerTCP()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("addr1=%v addr2=%v", addr1, addr2)
}
