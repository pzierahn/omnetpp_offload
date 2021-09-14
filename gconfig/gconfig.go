package gconfig

import (
	"fmt"
)

const (
	defaultBrokerPort   = 8888
	defaultStargatePort = 8889
)

type Configfile struct {
	Broker struct {
		Address      string `json:"address"`
		BrokerPort   int    `json:"port"`
		StargatePort int    `json:"stargatePort"`
	} `json:"broker"`
	Worker struct {
		Name string `json:"name"`
		// TODO: Rename to jobs
		Jobs int `json:"jobs"`
	} `json:"provider"`
}

func BrokerPort() (port int) {
	return Config.Broker.BrokerPort
}

func StargateAddr() (addr string) {
	addr = fmt.Sprintf("%s:%d", Config.Broker.Address, Config.Broker.StargatePort)
	return
}

func Jobs() (cpus int) {
	return Config.Worker.Jobs
}

func StargatePort() (port int) {
	return Config.Broker.StargatePort
}

func BrokerAddr() (addr string) {
	return Config.Broker.Address
}

func StargateDialAddr() (addr string) {
	addr = fmt.Sprintf("%s:%d", Config.Broker.Address, Config.Broker.StargatePort)
	return
}

func BrokerDialAddr() (addr string) {
	addr = fmt.Sprintf("%s:%d", Config.Broker.Address, Config.Broker.BrokerPort)
	return
}
