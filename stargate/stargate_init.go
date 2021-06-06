package stargate

import (
	"log"
	"net"
)

const (
	stunAddr = "31.18.129.212"
	stunPort = 9595
)

var rendezvousAddr = &net.UDPAddr{
	IP:   net.ParseIP(stunAddr),
	Port: stunPort,
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Stargate ")
}
