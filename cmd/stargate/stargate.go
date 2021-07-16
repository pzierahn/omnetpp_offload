package main

import (
	"context"
	"flag"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var timeout time.Duration
	var write string

	flag.StringVar(&write, "write", "", "")
	flag.DurationVar(&timeout, "timeout", time.Minute*8, "")
	flag.Parse()

	ctx, cnl := context.WithTimeout(context.Background(), timeout)
	defer cnl()

	conn, peer, err := stargate.DialUDP(ctx, "123456")
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
