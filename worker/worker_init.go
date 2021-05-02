package worker

import (
	"github.com/patrickz98/project.go.omnetpp/defines"
	"log"
	"os"
	"path/filepath"
)

var cachePath string

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "Worker ", log.LstdFlags|log.Lshortfile)

	cachePath = filepath.Join(defines.CacheDir(), "simulations")
	_ = os.MkdirAll(cachePath, 0755)
}

func Clean() {
	logger.Printf("cleaning worker cache %s\n", cachePath)
	err := os.RemoveAll(cachePath)
	if err != nil {
		panic(err)
	}
}
