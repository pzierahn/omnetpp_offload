package gconfig

import (
	"fmt"
)

const (
	defaultBrokerPort = 8888
)

type Configfile struct {
	Broker struct {
		Address      string `json:"address"`
		BrokerPort   int    `json:"port"`
		StargatePort int    `json:"stargatePort"`
	} `json:"broker"`
	Worker struct {
		Name string `json:"name"`
		Jobs int    `json:"jobs"`
	} `json:"provider"`
}

func BrokerPort() (port int) {
	return Config.Broker.BrokerPort
}

func Jobs() (jobs int) {
	return Config.Worker.Jobs
}

func StargatePort() (port int) {
	return Config.Broker.StargatePort
}

func BrokerAddr() (addr string) {
	return Config.Broker.Address
}

func BrokerDialAddr() (addr string) {
	addr = fmt.Sprintf("%s:%d", Config.Broker.Address, Config.Broker.BrokerPort)
	return
}
