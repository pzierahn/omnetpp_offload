package main

import (
	"encoding/json"
	"flag"
	"github.com/pzierahn/project.go.omnetpp/consumer"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"io/ioutil"
	"log"
	"path/filepath"
)

var path string
var configPath string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.StringVar(&path, "path", ".", "simulation path")
	flag.StringVar(&configPath, "config", "opp-edge-config.json", "simulation config JSON")
}

func main() {

	gconfig.ParseFlags()

	path, err := filepath.Abs(path)
	if err != nil {
		log.Panicln(err)
	}

	var runConfig consumer.Config
	runConfig.Path = path

	byt, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Panicln(err)
	}

	err = json.Unmarshal(byt, &runConfig)
	if err != nil {
		log.Panicln(err)
	}

	consumer.Start(&runConfig)
}
