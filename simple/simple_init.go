package simple

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Simple ", log.LstdFlags|log.Lshortfile)
}
