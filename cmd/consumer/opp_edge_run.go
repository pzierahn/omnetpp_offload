package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/pzierahn/project.go.omnetpp/consumer"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

var path string
var configPath string
var deadline time.Duration

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.StringVar(&path, "path", ".", "set simulation path")
	flag.StringVar(&configPath, "config", "opp-edge-config.json", "set simulation config json")
	flag.DurationVar(&deadline, "timeout", time.Hour*3, "set timeout for execution")
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

	ctx, cnl := context.WithTimeout(context.Background(), deadline)
	defer cnl()

	consumer.Start(ctx, &runConfig)
}
