package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/pzierahn/project.go.omnetpp/consumer"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

var path string
var configPath string
var timeout time.Duration

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.StringVar(&path, "path", ".", "set simulation path")
	flag.StringVar(&configPath, "config", "opp-edge-config.json", "set simulation config json")
	flag.DurationVar(&timeout, "timeout", time.Hour*3, "set timeout for execution")
}

func main() {

	gconfig.ParseFlags(gconfig.ParseBroker)

	stargate.SetConfig(stargate.Config{
		Addr: gconfig.BrokerAddr(),
		Port: gconfig.StargatePort(),
	})

	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalln(err)
	}

	var runConfig consumer.Config
	runConfig.Path = path

	byt, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(byt, &runConfig)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cnl := context.WithTimeout(context.Background(), timeout)
	defer cnl()

	consumer.OffloadSimulation(ctx, &runConfig)
}
