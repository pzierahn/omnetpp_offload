package main

import (
	"com.github.patrickz98.omnet/defines"
	"com.github.patrickz98.omnet/simple"
	"com.github.patrickz98.omnet/worker"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
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

func configName(deviceName string) {
	cleaner := regexp.MustCompile(`[^a-zA-B0-9-]`)

	idName := deviceName
	idName = strings.ToLower(idName)
	idName = cleaner.ReplaceAllString(idName, "_")
	idName = strings.Trim(idName, "_ ")

	config.WorkerId = idName + "-" + simple.RandomId(6)
	config.DeviceName = deviceName
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
		configName(deviceName)
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

	if err := worker.Link(config); err != nil {
		panic(err)
	}
}
