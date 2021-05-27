package gconfig

import "net"

type GRPCConnection struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (config GRPCConnection) DialAddr() (addr string) {
	//addr = fmt.Sprintf("[%s]::%d", config.Address, config.Port)
	udp := net.UDPAddr{
		IP:   net.ParseIP(config.Address),
		Port: config.Port,
	}

	addr = udp.String()

	return
}
