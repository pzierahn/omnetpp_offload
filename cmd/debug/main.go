package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx1, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx2, cancel1 := context.WithTimeout(ctx1, time.Second*5)
	defer cancel1()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		<-ctx1.Done()
		log.Printf("ctx1 Done")
	}()
	go func() {
		defer wg.Done()
		<-ctx2.Done()
		log.Printf("ctx2 Done")
	}()

	wg.Wait()

	//select {
	//case <-ctx1.Done():
	//	log.Printf("ctx1 Done")
	//case <-ctx2.Done():
	//	log.Printf("ctx2 Done")
	//}

	//addr1, addr2, err := stargate.RelayServerTCP()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Printf("addr1=%v addr2=%v", addr1, addr2)
}
