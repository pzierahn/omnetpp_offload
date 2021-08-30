package consumer

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/equic"
	"github.com/pzierahn/project.go.omnetpp/eval"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type providerConnection struct {
	ctx      context.Context
	info     *pb.ProviderInfo
	provider pb.ProviderClient
	store    pb.StorageClient
}

func (pConn *providerConnection) name() (name string) {
	return pConn.info.ProviderId
}

func (cons *consumer) connect(prov *pb.ProviderInfo) (conn *grpc.ClientConn, err error) {

	conn, err = cons.connectLocal(prov)
	if err == nil {
		eval.LogSetup(eval.ConnectLocal, prov)
		return
	}

	conn, err = cons.connectP2P(prov)
	if err == nil {
		eval.LogSetup(eval.ConnectP2P, prov)
		return
	}

	conn, err = cons.connectRelay(prov)

	if err == nil {
		eval.LogSetup(eval.ConnectRelay, prov)
	}

	return
}

func (cons *consumer) connectP2P(prov *pb.ProviderInfo) (cc *grpc.ClientConn, err error) {

	log.Printf("connectP2P: %v", prov.ProviderId)

	ctx, cln := context.WithTimeout(context.Background(), time.Second*5)
	defer cln()

	return equic.P2PDialGRPC(ctx, prov.ProviderId)
}

func (cons *consumer) connectRelay(prov *pb.ProviderInfo) (cc *grpc.ClientConn, err error) {

	log.Printf("connectRelay: %v", prov.ProviderId)

	ctx, cln := context.WithTimeout(cons.ctx, time.Second*5)
	defer cln()

	conn, err := stargate.DialRelayTCP(ctx, prov.ProviderId)
	if err != nil {
		return
	}

	log.Printf("connectRelay: dial %v", conn.RemoteAddr().String())

	return grpc.DialContext(
		ctx,
		conn.RemoteAddr().String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return conn, nil
		}),
	)
}

func (cons *consumer) connectLocal(prov *pb.ProviderInfo) (cc *grpc.ClientConn, err error) {

	log.Printf("connectLocal: %v", prov.ProviderId)

	ctx, cln := context.WithTimeout(cons.ctx, time.Second)
	defer cln()

	addr, err := stargate.DialLocal(ctx, prov.ProviderId)
	if err != nil {
		return
	}

	log.Printf("connectLocal: dial %v", addr)

	return grpc.DialContext(
		ctx,
		addr.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}
