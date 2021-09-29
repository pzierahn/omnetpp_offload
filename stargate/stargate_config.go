package stargate

import (
	"fmt"
	"net"
)

type Config struct {
	Addr string
	Port int
}

func (c *Config) tcpAddr() (addr *net.TCPAddr, err error) {
	// TODO: Check for ipv6.
	addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.Addr, c.Port))
	return
}

func (c *Config) udpAddr() (addr *net.UDPAddr, err error) {
	// TODO: Check for ipv6.
	addr, err = net.ResolveUDPAddr("tcp", fmt.Sprintf("%s:%d", c.Addr, c.Port))
	return
}
