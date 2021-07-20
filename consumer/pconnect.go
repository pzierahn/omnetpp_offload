package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
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

	conn, err := stargate.RelayDialTCP(ctx, prov.ProviderId)
	if err != nil {
		return
	}

	log.Printf("connectRelay: dial %v", conn.RemoteAddr().String())

	return grpc.Dial(
		conn.RemoteAddr().String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}
