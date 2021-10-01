package stargate

import (
	"context"
	"net"
	"time"
)

// DialP2PUDP will return an UDP connection and the peers UDP address.
// The connection is already established and tested.
func DialP2PUDP(ctx context.Context, dialAddr DialAddr) (conn *net.UDPConn, addr *net.UDPAddr, err error) {

	log.Printf("DialP2PUDP: dialAddr=%v", dialAddr)

	conn, err = net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}

	pr := peerResolver{
		conn: conn,
		dial: dialAddr,
	}
	peer, err := pr.resolvePeer(ctx)
	if err != nil {
		log.Printf("############# error: %v", err)
		return
	}

	addr = peer.Address

	log.Printf("DialP2PUDP: resolved peer dialAddr=%v peer=%v", dialAddr, peer)

	helper := p2pConnector{
		conn:    conn,
		start:   peer.Peer == 0,
		peer:    addr,
		timeout: time.Second * 2,
	}

	err = helper.connect(ctx)
	if err != nil {
		return
	}

	log.Printf("DialP2PUDP: dialAddr=%v peer=%v connection established", dialAddr, peer)

	return
}
