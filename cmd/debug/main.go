package main

import (
	"log"
	"sync"
	"time"
)

func listen(name string, a map[string]int, c *sync.Cond) {
	log.Println(name, " Lock")
	c.L.Lock()
	log.Println(name, " Wait")
	c.Wait()
	log.Println(name, " age:", a["T"])
	c.L.Unlock()
}

func broadcast(name string, a map[string]int, c *sync.Cond) {
	time.Sleep(time.Second)
	log.Println(name, " Lock")
	c.L.Lock()
	a["T"] = 25
	log.Println(name, " Broadcast")
	c.Broadcast()
	c.L.Unlock()
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cond := sync.NewCond(&sync.Mutex{})

	var val int

	go func() {
		for {
			log.Println("###1 Lock")
			cond.L.Lock()
			log.Println("###1 Wait")
			cond.Wait()
			log.Println("###1 val", val)
			log.Println("###1 Unlock")
			cond.L.Unlock()
		}
	}()

	//go func() {
	//	for {
	//		log.Println("###2 Lock")
	//		cond.L.Lock()
	//		log.Println("###2 Wait")
	//		cond.Wait()
	//		log.Println("###2 val", val)
	//		log.Println("###2 Unlock")
	//		cond.L.Unlock()
	//	}
	//}()
	//
	//go func() {
	//	for range time.Tick(time.Second) {
	//		cond.L.Lock()
	//		val++
	//		cond.Broadcast()
	//		cond.L.Unlock()
	//	}
	//}()

	cond.L.Lock()
	log.Println("val test init", val)
	cond.L.Unlock()

	for range time.Tick(time.Second) {
		log.Println("###2 Lock")
		cond.L.Lock()
		val++
		log.Println("###2 Broadcast")
		cond.Broadcast()
		log.Println("###2 Unlock")
		cond.L.Unlock()
	}

	//var age = make(map[string]int)
	//
	//var mu sync.Mutex
	//cond := sync.NewCond(&mu)
	//
	//go func() {
	//	for inx := 0; inx < 5; inx++ {
	//		listen("lis1", age, cond)
	//	}
	//}()
	//
	////// listener 1
	////go listen("lis1", age, cond)
	//
	////// listener 2
	////go listen("lis2", age, cond)
	//
	//// broadcast
	//go broadcast("b1", age, cond)
	//go broadcast("b2", age, cond)
	//go broadcast("b3", age, cond)
	//
	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, os.Interrupt)
	//<-ch
}
