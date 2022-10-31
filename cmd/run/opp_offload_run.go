package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/pzierahn/omnetpp_offload/consumer"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	"github.com/pzierahn/omnetpp_offload/simple"
	"log"
	"os"
	"path/filepath"
	"time"
)

var path string
var configPath string
var timeout time.Duration
var writeLog bool

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.StringVar(&path, "path", ".", "set simulation path")
	flag.StringVar(&configPath, "config", "", "set simulation config JSON")
	flag.DurationVar(&timeout, "timeout", time.Hour*3, "set timeout for execution")
	flag.BoolVar(&writeLog, "wl", false, "write logs to .cache/evaluation")

	flag.Parse()
}

func main() {

	if writeLog {
		simple.WriteLogToFile("consumer", gconfig.CacheDir())
	}

	config := gconfig.ParseFlagsBroker()

	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalln(err)
	}

	if configPath == "" {
		configPath = filepath.Join(path, "opp-offload-config.json")
	}

	var runConfig consumer.Config
	runConfig.Path = path
	runConfig.Scenario = os.Getenv("SCENARIO")
	runConfig.Trail = os.Getenv("TRAIL")

	byt, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(byt, &runConfig)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cnl := context.WithTimeout(context.Background(), timeout)
	defer cnl()

	consumer.OffloadSimulation(ctx, config, &runConfig)
}
