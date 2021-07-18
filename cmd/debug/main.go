package main

import (
	"context"
	"log"
	"os/exec"
	"time"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "sleep", "5").Run(); err != nil {
		// This will fail after 100 milliseconds. The 5 second sleep
		// will be interrupted.
	}

	//addr1, addr2, err := stargate.RelayServerTCP()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Printf("addr1=%v addr2=%v", addr1, addr2)
}
