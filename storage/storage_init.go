package storage

import (
	"com.github.patrickz98.omnet/defines"
	"log"
	"os"
)

const (
	storageAddress = "192.168.0.11:50052"
	storagePath    = defines.DataPath + "/storage"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Storage ", log.LstdFlags|log.Lshortfile)
}
