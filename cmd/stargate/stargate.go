package main

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"log"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(inx int) {
		defer wg.Done()

		//ctx, cnl := context.WithTimeout(context.Background(), time.Second*5)
		//defer cnl()
		ctx := context.Background()

		conn, peer, err := stargate.DialUDP(ctx, "123456")
		if err != nil {
			log.Printf("error %d: %v", inx, err)
			return
		}

		log.Printf("Connected %d: local=%v peer=%v", inx, conn.LocalAddr(), peer)

		wait := time.Second * 10
		time.Sleep(wait)
		log.Printf("Writing stuff...")

		_, err = conn.WriteTo([]byte("Pups"), peer)
		if err != nil {
			log.Printf("error %d: %v", inx, err)
		}

	}(1)

	wait := time.Second * 1
	log.Printf("Start waiting for %v...", wait)
	time.Sleep(wait)
	log.Printf("Done waiting")

	wg.Add(1)
	go func(inx int) {
		defer wg.Done()

		ctx, cnl := context.WithTimeout(context.Background(), time.Second*5)
		defer cnl()

		conn, peer, err := stargate.DialUDP(ctx, "123456")
		if err != nil {
			log.Printf("error %d: %v", inx, err)
			return
		}

		log.Printf("Connected %d: local=%v peer=%v", inx, conn.LocalAddr(), peer)

		buf := make([]byte, 1024)
		br, err := conn.Read(buf)
		if err != nil {
			log.Printf("error %d: %v", inx, err)
		}

		log.Printf("Connected %d: read %s", inx, string(buf[:br]))
	}(2)

	wg.Wait()
}
