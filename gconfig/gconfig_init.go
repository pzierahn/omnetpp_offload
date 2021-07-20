package gconfig

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/defines"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var Config Configfile

func init() {

	defer func() {
		//
		// Set default values
		//

		if Config.Broker.BrokerPort == 0 {
			Config.Broker.BrokerPort = defaultBrokerPort
		}

		if Config.Broker.StargatePort == 0 {
			Config.Broker.StargatePort = defaultStargatePort
		}

		if Config.Worker.Name == "" {
			Config.Worker.Name = runtime.GOOS + "-" + runtime.GOARCH
		}

		if Config.Worker.DevoteCPUs == 0 {
			Config.Worker.DevoteCPUs = runtime.NumCPU()
		}
	}()

	configPath := defines.ConfigDir()
	configFile := filepath.Join(configPath, "configuration.json")

	if _, err := os.Stat(configFile); err != nil {
		return
	}

	byt, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byt, &Config)
	if err != nil {
		panic(err)
	}
}

func ParseFlags() {

	//
	// Broker
	//
	flag.StringVar(&Config.Broker.Address, "broker", Config.Broker.Address, "set broker address")
	flag.IntVar(&Config.Broker.BrokerPort, "port", Config.Broker.BrokerPort, "set broker port")
	flag.IntVar(&Config.Broker.StargatePort, "stargatePort", Config.Broker.StargatePort, "set stargate port")

	//
	// Worker
	//
	flag.StringVar(&Config.Worker.Name, "name", Config.Worker.Name, "set worker name")
	flag.IntVar(&Config.Worker.DevoteCPUs, "devoteCPUs", Config.Worker.DevoteCPUs, "set how manny CPUs should be used")

	flag.Parse()
}

func Persist() {
	configPath := defines.ConfigDir()

	err := os.MkdirAll(configPath, 0755)
	if err != nil {
		panic(err)
	}

	byt, err := json.MarshalIndent(Config, "", "  ")
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(configPath, "configuration.json")
	fmt.Println("write Config to", configFile)

	err = ioutil.WriteFile(configFile, byt, 0644)
	if err != nil {
		panic(err)
	}
}
