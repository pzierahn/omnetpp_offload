package gconfig

import (
	"encoding/json"
	"flag"
	"github.com/pzierahn/omnetpp_offload/simple"
	"github.com/pzierahn/omnetpp_offload/stargate"
	"os"
	"path/filepath"
	"runtime"
)

const (
	parseBroker = 1 << iota
	parseWorker
	parseAll = parseBroker | parseWorker
)

var defaultConfig Config

func init() {

	defer func() {
		//
		// Set default values
		//

		if defaultConfig.Broker.BrokerPort == 0 {
			defaultConfig.Broker.BrokerPort = DefaultBrokerPort
		}

		if defaultConfig.Broker.StargatePort == 0 {
			defaultConfig.Broker.StargatePort = stargate.DefaultPort
		}

		if defaultConfig.Provider.Name == "" {
			//Config.Provider.Name = runtime.GOOS + "-" + runtime.GOARCH
			defaultConfig.Provider.Name = simple.GetHostnameShort()
		}

		if defaultConfig.Provider.Jobs == 0 {
			defaultConfig.Provider.Jobs = runtime.NumCPU()
		}
	}()

	configPath := ConfigDir()
	configFile := filepath.Join(configPath, "configuration.json")

	if _, err := os.Stat(configFile); err != nil {
		return
	}

	byt, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byt, &defaultConfig)
	if err != nil {
		panic(err)
	}
}

func parseFlags(parse int) {

	if parseBroker&parse != 0 {
		//
		// Broker command line arguments
		//
		flag.StringVar(&defaultConfig.Broker.Address, "broker", defaultConfig.Broker.Address, "set broker address")
		flag.IntVar(&defaultConfig.Broker.BrokerPort, "port", defaultConfig.Broker.BrokerPort, "set broker port")
		flag.IntVar(&defaultConfig.Broker.StargatePort, "stargate", defaultConfig.Broker.StargatePort, "set stargate port")
	}

	if parseWorker&parse != 0 {
		//
		// Worker command line arguments
		//
		flag.StringVar(&defaultConfig.Provider.Name, "name", defaultConfig.Provider.Name, "set worker name")
		flag.IntVar(&defaultConfig.Provider.Jobs, "jobs", defaultConfig.Provider.Jobs, "set how many jobs should be started")
	}

	flag.Parse()
}

func Default() Config {
	return defaultConfig
}

func ParseFlags() Config {
	parseFlags(parseAll)
	return defaultConfig
}

func ParseFlagsBroker() Broker {
	parseFlags(parseBroker)
	return defaultConfig.Broker
}
