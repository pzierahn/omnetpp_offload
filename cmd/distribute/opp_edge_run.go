package main

import (
	"encoding/json"
	"flag"
	"github.com/patrickz98/project.go.omnetpp/distribute"
	"io/ioutil"
	"log"
	"path/filepath"
)

var path string
var configPath string

func init() {
	flag.StringVar(&path, "path", ".", "simulation path")
	flag.StringVar(&configPath, "config", "opp-edge-config.json", "simulation config JSON")
}

func main() {

	flag.Parse()

	path, err := filepath.Abs(path)
	if err != nil {
		log.Panicln(err)
	}

	var config distribute.Config
	config.Path = path

	byt, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Panicln(err)
	}

	err = json.Unmarshal(byt, &config)
	if err != nil {
		log.Panicln(err)
	}

	err = distribute.Run(&config)
	if err != nil {
		panic(err)
	}
}
