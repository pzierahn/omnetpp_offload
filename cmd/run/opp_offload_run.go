package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/consumer"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	"io"
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

func enableLogWrite() {

	date := time.Now()
	logName := fmt.Sprintf("consumer.%s.%s.log",
		date.Format("2006-02-01"), date.Format("15-04-05"))

	logDir := filepath.Join(gconfig.CacheDir(), "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalln(err)
	}

	logPath := filepath.Join(logDir, logName)
	log.Printf("Write logs to %s", logPath)

	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	mw := io.MultiWriter(os.Stderr, file)
	log.SetOutput(mw)
}

func main() {

	if writeLog {
		enableLogWrite()
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
