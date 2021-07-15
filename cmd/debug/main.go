package main

import (
	"context"
	"log"
	"time"
)

//func test() {
//	ch := make(chan bool)
//	defer func() {
//		log.Println("Closing ch")
//		close(ch)
//	}()
//
//	go func() {
//		time.Sleep(time.Second * 1)
//		log.Println("Writing to ch")
//		ch <- true
//	}()
//}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)

	deadline, ok := ctx.Deadline()
	log.Printf("deadline=%v ok=%v", deadline, ok)

	//test()
	//
	//time.Sleep(time.Second * 3)
	//log.Println("Done")
}
