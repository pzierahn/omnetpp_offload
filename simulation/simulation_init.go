package simulation

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Simulation ", log.LstdFlags|log.Lshortfile)
}
