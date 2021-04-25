package storage

import (
	"github.com/patrickz98/project.go.omnetpp/defines"
	"log"
	"os"
	"path/filepath"
)

const (
	storageAddress = "192.168.0.11:50051"
)

var storagePath = filepath.Join(defines.DataPath, "/Storage")

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Storage ", log.LstdFlags|log.Lshortfile)
}
