package main

import (
	"encoding/json"
	"flag"
	"github.com/patrickz98/project.go.omnetpp/distribute"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	"io/ioutil"
	"log"
	"path/filepath"
)

var path string
var configPath string
var status bool

var config gconfig.Config

func init() {
	flag.StringVar(&path, "path", ".", "simulation path")
	flag.StringVar(&configPath, "config", "opp-edge-config.json", "simulation config JSON")
	flag.BoolVar(&status, "status", false, "status for simulationIds")

	config = gconfig.SourceAndParse(gconfig.ParseBroker)
}

func main() {

	if status {
		distribute.Status(config.Broker, nil)
		return
	}

	path, err := filepath.Abs(path)
	if err != nil {
		log.Panicln(err)
	}

	var runConfig distribute.Config
	runConfig.Path = path

	byt, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Panicln(err)
	}

	err = json.Unmarshal(byt, &runConfig)
	if err != nil {
		log.Panicln(err)
	}

	err = distribute.Run(config.Broker, &runConfig)
	if err != nil {
		panic(err)
	}
}
