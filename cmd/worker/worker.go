package main

import (
	"com.github.patrickz98.omnet/defines"
	"com.github.patrickz98.omnet/simple"
	"com.github.patrickz98.omnet/worker"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

const configPath = defines.DataPath + "/worker-config.json"

var deviceName string
var brokerAddress string

var configure bool

var config worker.Config

func init() {
	flag.StringVar(&deviceName, "deviceName", "", "set workerId")
	flag.StringVar(&brokerAddress, "brokerAddress", "", "set broker server address")

	flag.BoolVar(&configure, "configure", false, "generate config file with params")

	if _, err := os.Stat(configPath); err == nil {
		_ = simple.UnmarshallFile(configPath, &config)
	}
}

func persistConfig() {

	byt, _ := json.MarshalIndent(config, "", "  ")

	fmt.Println("config", string(byt))
	fmt.Println("write config to", configPath)

	err := ioutil.WriteFile(configPath, byt, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {

	flag.Parse()

	if deviceName != "" {
		config.WorkerId = simple.NamedId(deviceName, 6)
		config.DeviceName = deviceName
	}

	if brokerAddress != "" {
		config.BrokerAddress = brokerAddress
	}

	if configure {

		//
		// Persist config
		//

		persistConfig()
		return
	}

	if config.WorkerId == "" {

		//
		// No name configured
		//

		config.WorkerId = fmt.Sprintf("%s-fish-%s",
			runtime.GOOS, simple.RandomId(6))
	}

	if config.BrokerAddress == "" {
		config.BrokerAddress = defines.Port
	}

	if brokerAddress != "" {
		config.BrokerAddress = brokerAddress
	}

	conn, err := worker.Connect(config)
	if err != nil {
		panic(err)
	}

	if err := conn.StartLink(context.Background()); err != nil {
		panic(err)
	}
}
