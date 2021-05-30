package provider

import (
	"github.com/patrickz98/project.go.omnetpp/defines"
	"log"
	"os"
	"path/filepath"
)

var cachePath string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cachePath = filepath.Join(defines.CacheDir(), "simulations")
	_ = os.MkdirAll(cachePath, 0755)
}

func Clean() {
	log.Printf("cleaning worker cache %s\n", cachePath)
	err := os.RemoveAll(cachePath)
	if err != nil {
		log.Fatalln(err)
	}
}
