package gconfig

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/defines"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func Source() (config Config) {

	defer func() {
		//
		// Set default values
		//

		if config.Broker.Port == 0 {
			config.Broker.Port = defines.DefaultPort
		}

		if config.Worker.Name == "" {
			config.Worker.Name = simple.GetHostnameShort()
		}

		if config.Worker.DevoteCPUs == 0 {
			config.Worker.DevoteCPUs = runtime.NumCPU()
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

	err = json.Unmarshal(byt, &config)
	if err != nil {
		panic(err)
	}

	return
}

func SourceAndParse(parse ...int) (config Config) {
	config = Source()

	if len(parse) == 0 {
		parse = ParseAll
	}

	check := make(map[int]bool)

	for _, elem := range parse {
		check[elem] = true
	}

	if check[ParseBroker] {
		flag.StringVar(&config.Broker.Address, "broker", config.Broker.Address, "set broker address")
		flag.IntVar(&config.Broker.Port, "port", config.Broker.Port, "set broker port")
	}

	if check[ParseWorker] {
		flag.IntVar(&config.Worker.DevoteCPUs, "devoteCPUs", config.Worker.DevoteCPUs, "set how manny CPUs should be used")
		flag.StringVar(&config.Worker.Name, "name", config.Worker.Name, "name worker client")
	}

	flag.Parse()

	return
}

func Persist(config Config) {
	configPath := defines.ConfigDir()

	err := os.MkdirAll(configPath, 0755)
	if err != nil {
		panic(err)
	}

	byt, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(configPath, "configuration.json")
	fmt.Println("write config to", configFile)

	err = ioutil.WriteFile(configFile, byt, 0644)
	if err != nil {
		panic(err)
	}

	return
}
