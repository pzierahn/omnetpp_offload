package storage

import (
	"github.com/pzierahn/project.go.omnetpp/defines"
	lg "log"
	"os"
	"path/filepath"
)

var storagePath string

var log *lg.Logger

func init() {
	log = lg.New(os.Stderr, "Storage ", lg.LstdFlags|lg.Lshortfile)

	storagePath = filepath.Join(defines.CacheDir(), "storage")
	_ = os.MkdirAll(storagePath, 0755)
}

func Clean() {
	log.Printf("Clean: %v\n", storagePath)
	_ = os.RemoveAll(storagePath)
}
