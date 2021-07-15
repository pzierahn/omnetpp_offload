package main

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"log"
	"sync"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(inx int) {
		defer wg.Done()

		conn, remote, err := stargate.DialUDP(context.Background(), "123456")
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Connect %d: local=%v remote=%v", inx, conn.LocalAddr(), remote)
	}(1)

	wg.Add(1)
	go func(inx int) {
		defer wg.Done()

		conn, remote, err := stargate.DialUDP(context.Background(), "123456")
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Connect %d: local=%v remote=%v", inx, conn.LocalAddr(), remote)
	}(2)

	wg.Wait()
}
