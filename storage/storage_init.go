package storage

import (
	"github.com/pzierahn/project.go.omnetpp/defines"
	"log"
	"os"
	"path/filepath"
)

var storagePath string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	storagePath = filepath.Join(defines.CacheDir(), "storage")
	_ = os.MkdirAll(storagePath, 0755)
}

func Clean() {
	log.Printf("cleaning storage %s\n", storagePath)
	err := os.RemoveAll(storagePath)
	if err != nil {
		log.Fatalln(err)
	}
}
