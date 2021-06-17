package stargate

import (
	lg "log"
	"net"
	"os"
)

const (
	stunAddr = "31.18.129.212"
	stunPort = 9595
)

var rendezvousAddr = &net.UDPAddr{
	IP:   net.ParseIP(stunAddr),
	Port: stunPort,
}

var log *lg.Logger

func init() {
	log = lg.New(os.Stderr, "Stargate ", lg.LstdFlags|lg.Lshortfile)
}
