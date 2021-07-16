package main

import (
	"context"
	"flag"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"log"
	"net"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var server bool
	var timeout time.Duration
	var dialAddr string
	var write string

	flag.BoolVar(&server, "server", false, "start stun server")
	flag.StringVar(&dialAddr, "dialAddr", "", "dial address")
	flag.StringVar(&write, "write", "", "the message that will be transferred")
	flag.DurationVar(&timeout, "timeout", time.Minute*8, "timeout for connection")
	flag.Parse()

	if server {
		err := stargate.Server(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}

	if dialAddr == "" {
		log.Fatalln("dialAddr missing!")
	}

	// Set stun server address
	stargate.SetRendezvousServer(&net.UDPAddr{
		IP:   net.ParseIP("31.18.129.212"),
		Port: 9595,
	})

	ctx, cnl := context.WithTimeout(context.Background(), timeout)
	defer cnl()

	conn, peer, err := stargate.DialUDP(ctx, dialAddr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Connected: local=%v peer=%v", conn.LocalAddr(), peer)

	if write != "" {
		log.Printf("Write: '%s'", write)

		_, err = conn.WriteTo([]byte(write), peer)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		buf := make([]byte, 1024)
		br, err := conn.Read(buf)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Read: '%s'", string(buf[:br]))
	}
}
