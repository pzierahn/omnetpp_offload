package equic

import (
	"context"
	"github.com/lucas-clemente/quic-go"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"net"
)

func P2PDialGRPC(ctx context.Context, dialAddr stargate.DialAddr) (conn *grpc.ClientConn, err error) {
	gate, remote, err := stargate.DialP2PUDP(ctx, dialAddr)
	if err != nil {
		return
	}

	return grpc.Dial(
		remote.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(dialAdapter(gate)),
	)
}

func P2PListenerQUIC(ctx context.Context, dialAddr stargate.DialAddr) (p2p quic.Listener, err error) {

	conn, _, err := stargate.DialP2PUDP(ctx, dialAddr)
	if err != nil {
		return
	}

	tlsConf, _ := generateTLSConfig()

	return quic.Listen(conn, tlsConf, config)
}

func P2PListener(ctx context.Context, dialAddr stargate.DialAddr) (p2p net.Listener, err error) {

	qLis, err := P2PListenerQUIC(ctx, dialAddr)
	if err != nil {
		return
	}

	return Listen(qLis), err
}
