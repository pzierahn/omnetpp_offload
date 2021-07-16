package stargate

import (
	lg "log"
	"net"
	"os"
)

const (
	defaultAddr = "31.18.129.212"
	defaultPort = 9595
)

var rendezvousAddr = &net.UDPAddr{
	IP:   net.ParseIP(defaultAddr),
	Port: defaultPort,
}

func SetRendezvousServer(addr *net.UDPAddr) {
	rendezvousAddr = addr
}

var log *lg.Logger

func init() {
	log = lg.New(os.Stderr, "Stargate ", lg.LstdFlags|lg.Lshortfile)
}
