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
		log.Println(err)

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

	ctx, cln := context.WithTimeout(context.Background(), time.Second*2)
	defer cln()

	gate := pb.NewStargateClient(cons.bconn)
	port, err := gate.Relay(ctx, &pb.RelayRequest{
		DialAddr: prov.ProviderId,
	})
	if err != nil {
		return
	}

	log.Printf("Connect over relay server (port: %v)", port.Port)

	// TODO: replace gconfig.Config.Broker.Address
	raddr := &net.TCPAddr{
		IP:   net.ParseIP(gconfig.Config.Broker.Address),
		Port: int(port.Port),
	}

	return grpc.Dial(
		raddr.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}
