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

	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	var value int

	var wg sync.WaitGroup

	go func() {
		for inx := 0; inx < 2; inx++ {

			log.Println("reader", inx, "lock")
			cond.L.Lock()
			log.Println("reader", inx, "lock wait")
			cond.Wait()
			log.Println("reader", inx, "value", value)
			cond.L.Unlock()
			log.Println("reader", inx, "done")
		}
	}()

	time.Sleep(time.Second * 2)

	for inx := 0; inx < 2; inx++ {
		wg.Add(1)

		go func(inx int) {
			defer wg.Done()

			log.Println("write", inx, "Waiting")
			cond.L.Lock()

			log.Println("write", inx, "doing stuff")
			time.Sleep(time.Second * 4)

			value = inx
			cond.Broadcast()

			cond.L.Unlock()
			log.Println("write", inx, "Done")
		}(inx)
	}

	wg.Wait()

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
