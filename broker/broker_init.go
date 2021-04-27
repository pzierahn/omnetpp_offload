package broker

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Broker ", log.LstdFlags|log.Lshortfile)
}
