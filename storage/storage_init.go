package storage

import (
	"github.com/pzierahn/project.go.omnetpp/defines"
	"log"
	"os"
	"path/filepath"
)

var storagePath string

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Storage ", log.LstdFlags|log.Lshortfile)

	storagePath = filepath.Join(defines.CacheDir(), "storage")
	_ = os.MkdirAll(storagePath, 0755)
}

func Clean() {
	logger.Printf("cleaning storage %s\n", storagePath)
	err := os.RemoveAll(storagePath)
	if err != nil {
		panic(err)
	}
}
