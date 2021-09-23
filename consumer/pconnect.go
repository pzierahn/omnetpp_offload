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
	connectAll = connectLocal | connectP2P | connectRelay
)

type providerConnection struct {
	ctx          context.Context
	conn         *grpc.ClientConn
	info         *pb.ProviderInfo
	provider     pb.ProviderClient
	store        pb.StorageClient
	downloadPipe chan *download
}

type download struct {
	task *pb.SimulationRun
	ref  *pb.StorageRef
}

func (pConn *providerConnection) id() (name string) {
	return pConn.info.ProviderId
}

func (pConn *providerConnection) close() {
	//TODO: pConn.provider.DropSession(ctx, &pb.Session{})

	close(pConn.downloadPipe)
	_ = pConn.conn.Close()
}

func pconnect(ctx context.Context, prov *pb.ProviderInfo) (conn *grpc.ClientConn, err error) {

	connect := connectAll

	// Eval stuff to ensure that only the desired connection will be used
	switch os.Getenv("CONNECT") {
	case "local":
		log.Println("########################## eval debug: connect only local!")
		connect = connectLocal

	case "p2p":
		log.Println("########################## eval debug: connect only p2p!")
		connect = connectP2P

	case "relay":
		log.Println("########################## eval debug: connect only relay!")
		connect = connectRelay
	}

	if connect&connectLocal != 0 {
		conn, err = pconnectLocal(ctx, prov.ProviderId)
		if err == nil {
			eval.LogSetup(eval.ConnectLocal, prov)
			return
		}
	}

	if connect&connectP2P != 0 {
		conn, err = pconnectP2P(ctx, prov.ProviderId)
		if err == nil {
			eval.LogSetup(eval.ConnectP2P, prov)
			return
		}
	}

	if connect&connectRelay != 0 {
		conn, err = pconnectRelay(ctx, prov.ProviderId)
		if err == nil {
			eval.LogSetup(eval.ConnectRelay, prov)
			return
		}
	}

	return
}

func pconnectP2P(ctx context.Context, providerId string) (cc *grpc.ClientConn, err error) {

	log.Printf("connectP2P: %v", providerId)

	ctx, cln := context.WithTimeout(ctx, time.Second*5)
	defer cln()

	return equic.P2PDialGRPC(ctx, providerId)
}

func pconnectRelay(ctx context.Context, providerId string) (cc *grpc.ClientConn, err error) {

	log.Printf("connectRelay: %v", providerId)

	ctx, cln := context.WithTimeout(ctx, time.Second*5)
	defer cln()

	conn, err := stargate.DialRelayTCP(ctx, providerId)
	if err != nil {
		return
	}

	log.Printf("connectRelay: dial %v", conn.RemoteAddr().String())

	return grpc.DialContext(
		ctx,
		conn.RemoteAddr().String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) {
			return conn, nil
		}),
	)
}

func pconnectLocal(ctx context.Context, providerId string) (cc *grpc.ClientConn, err error) {

	log.Printf("connectLocal: %v", providerId)

	ctx, cln := context.WithTimeout(ctx, time.Second)
	defer cln()

	addr, err := stargate.DialLocal(ctx, providerId)
	if err != nil {
		return
	}

	log.Printf("pconnectLocal: dial %v", addr)

	return grpc.DialContext(
		ctx,
		addr.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}
