package main

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"log"
	"time"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go func() {
		err := stargate.ServerRelayTCP()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	time.Sleep(time.Second)

	go func() {
		ctx, cnl := context.WithTimeout(context.Background(), time.Second*2)
		defer cnl()

		conn, err := stargate.RelayDialTCP(ctx, "1234567")
		if err != nil {
			log.Fatalln(err)
		}

		_, err = conn.Write([]byte("Hallo"))
		if err != nil {
			log.Fatalln(err)
		}
	}()

	//time.Sleep(time.Second * 6)

	conn, err := stargate.RelayDialTCP(context.Background(), "1234567")
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 1024)
	br, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Recieved: %v", string(buf[:br]))

	time.Sleep(time.Second * 10)
}
