package main

import (
	"context"
	"flag"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"log"
	"time"
)

var server bool
var timeout time.Duration
var dialAddr string
var write string
var serverAddr string
var port int

func init() {
	flag.BoolVar(&server, "server", false, "start stun server")
	flag.StringVar(&dialAddr, "dialAddr", "", "dial address")
	flag.StringVar(&write, "write", "", "the message that will be transferred")
	flag.DurationVar(&timeout, "timeout", time.Minute*8, "timeout for connection")
	flag.StringVar(&serverAddr, "addr", stargate.DefaultAddr, "set stargate server")
	flag.IntVar(&port, "port", stargate.DefaultPort, "set stargate server")
	flag.Parse()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set stun server address
	stargate.SetConfig(stargate.Config{
		Addr: serverAddr,
		Port: port,
	})

	if server {
		err := stargate.Server(context.Background(), true)
		if err != nil {
			log.Fatalln(err)
		}

		return
	}

	if dialAddr == "" {
		log.Fatalln("dialAddr missing!")
	}

	ctx, cnl := context.WithTimeout(context.Background(), timeout)
	defer cnl()

	conn, peer, err := stargate.DialP2PUDP(ctx, dialAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = conn.Close() }()

	log.Printf("Connected peer to peer: local=%v peer=%v", conn.LocalAddr(), peer)

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
		//log.Printf("Read: '%x'", buf[:br])
	}
}
