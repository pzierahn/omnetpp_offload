package broker

import "log"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Broker ")
}
