package provider

import (
	"github.com/pzierahn/project.go.omnetpp/defines"
	"log"
	"os"
	"path/filepath"
)

var cachePath string
var sessionsPath string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Provider ")

	cachePath = filepath.Join(defines.CacheDir(), "simulations")
	_ = os.MkdirAll(cachePath, 0755)

	sessionsPath = filepath.Join(defines.CacheDir(), "sessions.json")
}

func Clean() {
	log.Printf("Clean: %s\n", sessionsPath)
	_ = os.RemoveAll(sessionsPath)

	log.Printf("Clean: %s\n", cachePath)
	_ = os.RemoveAll(cachePath)
}
