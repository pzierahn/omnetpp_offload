package stateinfo

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Stateinfo ", log.LstdFlags|log.Lshortfile)
}
