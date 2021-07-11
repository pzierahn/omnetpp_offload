package consumer

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/equic"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type providerConnection struct {
	info     *pb.ProviderInfo
	conn     *net.UDPConn
	cConn    *grpc.ClientConn
	provider pb.ProviderClient
	store    pb.StorageClient
}

func (pConn *providerConnection) name() (name string) {
	return pConn.info.ProviderId
}

func connect(prov *pb.ProviderInfo) (conn *providerConnection, err error) {

	log.Printf("connect to provider %v", prov.ProviderId)

	ctx, cln := context.WithTimeout(context.Background(), time.Second*5)
	defer cln()

	gate, remote, err := stargate.Dial(ctx, prov.ProviderId)
	if err != nil {
		// Connection failed!
		return
	}

	var cConn *grpc.ClientConn
	cConn, err = grpc.Dial(
		remote.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(equic.GRPCDialer(gate)),
	)
	if err != nil {
		return
	}

	conn = &providerConnection{
		info:     prov,
		conn:     gate,
		cConn:    cConn,
		provider: pb.NewProviderClient(cConn),
		store:    pb.NewStorageClient(cConn),
	}

	return
}

func (pConn *providerConnection) close() {
	_ = pConn.cConn.Close()
	_ = pConn.conn.Close()
}
