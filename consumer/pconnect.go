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
	"os"
	"time"
)

const (
	connectLocal = 1 << iota
	connectP2P
	connectRelay
)

const connectAll = connectLocal | connectP2P | connectRelay

type providerConnection struct {
	ctx      context.Context
	info     *pb.ProviderInfo
	provider pb.ProviderClient
	store    pb.StorageClient
}

func (pConn *providerConnection) id() (name string) {
	return pConn.info.ProviderId
}

func (cons *consumer) connect(prov *pb.ProviderInfo) (conn *grpc.ClientConn, err error) {

	connect := connectAll

	// Eval stuff to ensure that only the desired connection will be established
	switch os.Getenv("CONNECT") {
	case "local":
		log.Println("########################## eval debug: connect only local!")
		connect = connectLocal

	case "p2p":
		log.Println("########################## eval debug: connect only p2p!")
		connect = connectP2P

	case "relay":
		log.Println("########################## eval debug: connect only local!")
		connect = connectRelay
	}

	if connect&connectLocal != 0 {
		conn, err = cons.connectLocal(prov)
		if err == nil {
			eval.LogSetup(eval.ConnectLocal, prov)
			return
		}
	}

	if connect&connectP2P != 0 {
		conn, err = cons.connectP2P(prov)
		if err == nil {
			eval.LogSetup(eval.ConnectP2P, prov)
			return
		}
	}

	if connect&connectRelay != 0 {
		conn, err = cons.connectRelay(prov)
		if err == nil {
			eval.LogSetup(eval.ConnectRelay, prov)
		}
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
