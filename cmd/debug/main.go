package main

import (
	"log"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	items := make([]int, 1)
	items[0] = 666
	item, items := items, items[1:]

	log.Println("item", item)
}
