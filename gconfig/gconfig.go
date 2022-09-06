package gconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultBrokerPort = 8888
)

type Broker struct {
	Address      string `json:"address"`
	BrokerPort   int    `json:"port"`
	StargatePort int    `json:"stargatePort"`
}

func (conf Broker) BrokerDialAddr() (addr string) {
	addr = fmt.Sprintf("%s:%d", conf.Address, conf.BrokerPort)
	return
}

type Provider struct {
	Name string `json:"name"`
	Jobs int    `json:"jobs"`
}

type Config struct {
	Broker   Broker   `json:"broker"`
	Provider Provider `json:"provider"`
}

func Write() {
	configPath := ConfigDir()

	err := os.MkdirAll(configPath, 0755)
	if err != nil {
		panic(err)
	}

	byt, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(configPath, "configuration.json")
	fmt.Println("write config to", configFile)

	err = os.WriteFile(configFile, byt, 0644)
	if err != nil {
		panic(err)
	}
}
