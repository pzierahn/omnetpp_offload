package common

import "fmt"

type GRPCConnection struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (config GRPCConnection) DialAddr() (addr string) {
	addr = fmt.Sprintf("%s:%d", config.Address, config.Port)
	return
}
