package stargate

import (
	"context"
	"github.com/lucas-clemente/quic-go"
	"github.com/pzierahn/project.go.omnetpp/equic"
	"google.golang.org/grpc"
	"net"
	"time"
)

func DialP2PUDP(ctx context.Context, dialAddr DialAddr) (conn *net.UDPConn, peer *net.UDPAddr, err error) {

	log.Printf("DialP2PUDP: dialAddr=%v", dialAddr)

	conn, err = net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}

	pr := peerResolver{
		conn: conn,
		dial: dialAddr,
	}
	peer, err = pr.resolvePeer(ctx)
	if err != nil {
		return
	}

	log.Printf("DialP2PUDP: dialAddr=%v peer=%v", dialAddr, peer)

	helper := p2pConnector{
		conn:      conn,
		peer:      peer,
		packages:  3,
		sendDelay: time.Millisecond * 20,
		timeout:   time.Millisecond * 600,
		// Use a sendDelay of 3s to test if both peers are in the same NAT
		//sendDelay: time.Second * 3,
		//timeout:   time.Second * 12,
	}

	err = helper.connect(ctx)
	if err != nil {
		return
	}

	// Everything is okay
	log.Printf("DialP2PUDP: dialAddr=%v peer=%v connection established", dialAddr, peer)

	return
}

func DialGRPC(ctx context.Context, dialAddr DialAddr) (conn *grpc.ClientConn, err error) {
	gate, remote, err := DialP2PUDP(ctx, dialAddr)
	if err != nil {
		// Connection failed!
		return
	}

	conn, err = grpc.Dial(
		remote.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(equic.GRPCDialer(gate)),
	)

	return
}

func ListenerQUIC(ctx context.Context, dialAddr DialAddr) (p2p quic.Listener, err error) {

	conn, _, err := DialP2PUDP(ctx, dialAddr)
	if err != nil {
		return
	}

	tlsConf, _ := equic.GenerateTLSConfig()

	return quic.Listen(conn, tlsConf, equic.Config)
}

func ListenerNet(ctx context.Context, dialAddr DialAddr) (p2p net.Listener, err error) {

	qLis, err := ListenerQUIC(ctx, dialAddr)
	if err != nil {
		return
	}

	return equic.Listen(qLis), err
}
