package main

import (
	"flag"
	"github.com/pzierahn/project.go.omnetpp/consumer"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"log"
	"path/filepath"
)

var path string
var configPath string

var config gconfig.Config

func init() {
	flag.StringVar(&path, "path", ".", "simulation path")
	flag.StringVar(&configPath, "config", "opp-edge-config.json", "simulation config JSON")

	config = gconfig.SourceAndParse(gconfig.ParseBroker)
}

func main() {

	path, err := filepath.Abs(path)
	if err != nil {
		log.Panicln(err)
	}

	var runConfig consumer.Config
	runConfig.Path = path

	//byt, err := ioutil.ReadFile(configPath)
	//if err != nil {
	//	log.Panicln(err)
	//}
	//
	//err = json.Unmarshal(byt, &runConfig)
	//if err != nil {
	//	log.Panicln(err)
	//}

	err = consumer.Run(config.Broker, &runConfig)
	if err != nil {
		panic(err)
	}
}
