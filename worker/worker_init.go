package worker

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Worker ", log.LstdFlags|log.Lshortfile)
}
