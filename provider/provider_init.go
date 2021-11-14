package provider

import (
	"github.com/pzierahn/omnetpp_offload/gconfig"
	"log"
	"os"
	"path/filepath"
)

var cachePath string
var sessionsPath string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Provider ")

	cachePath = filepath.Join(gconfig.CacheDir(), "simulations")
	_ = os.MkdirAll(cachePath, 0755)

	sessionsPath = filepath.Join(gconfig.CacheDir(), "sessions.json")
}

func Clean() {
	log.Printf("Clean: %s\n", sessionsPath)
	_ = os.RemoveAll(sessionsPath)

	log.Printf("Clean: %s\n", cachePath)
	_ = os.RemoveAll(cachePath)
}
