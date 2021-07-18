package stargate

import (
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	lg "log"
	"net"
	"os"
)

var rendezvousAddr *net.UDPAddr

func GetRendezvousServer() (addr *net.UDPAddr, err error) {

	if rendezvousAddr != nil {
		return rendezvousAddr, nil
	}

	return net.ResolveUDPAddr("udp", gconfig.StargateAddr())
}

func SetRendezvousServer(addr *net.UDPAddr) {
	rendezvousAddr = addr
}

var log *lg.Logger

func init() {
	log = lg.New(os.Stderr, "Stargate ", lg.LstdFlags|lg.Lshortfile)
}
