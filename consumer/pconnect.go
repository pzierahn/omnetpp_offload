package consumer

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
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

func (cons *consumer) connect(prov *pb.ProviderInfo) (conn *providerConnection, err error) {

	cc, err := cons.connectP2P(prov)
	if err != nil {
		log.Println(prov.ProviderId, err)

		cc, err = cons.connectRelay(prov)
		if err != nil {
			return
		}
	}

	conn = &providerConnection{
		info:     prov,
		provider: pb.NewProviderClient(cc),
		store:    pb.NewStorageClient(cc),
	}

	return
}

func (cons *consumer) connectP2P(prov *pb.ProviderInfo) (cc *grpc.ClientConn, err error) {

	log.Printf("connectP2P: %v", prov.ProviderId)

	ctx, cln := context.WithTimeout(context.Background(), time.Second*5)
	defer cln()

	return stargate.DialGRPCClientConn(ctx, prov.ProviderId)
}

func (cons *consumer) connectRelay(prov *pb.ProviderInfo) (cc *grpc.ClientConn, err error) {

	log.Printf("connectRelay: %v", prov.ProviderId)

	ctx, cln := context.WithTimeout(cons.ctx, time.Second*2)
	defer cln()

	gate := pb.NewStargateClient(cons.bconn)
	port, err := gate.Relay(ctx, &pb.RelayRequest{
		DialAddr: prov.ProviderId,
	})
	if err != nil {
		return
	}

	log.Printf("connectRelay: port=%v", port.Port)

	raddr, err := net.ResolveTCPAddr("tcp", gconfig.BrokerDialAddr())
	if err != nil {
		return
	}
	raddr.Port = int(port.Port)

	log.Printf("connectRelay: relay=%v", raddr.String())

	return grpc.Dial(
		raddr.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}
