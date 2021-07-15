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
	info     *pb.ProviderInfo
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

	cConn, err := stargate.DialGRPCClientConn(ctx, prov.ProviderId)
	if err != nil {
		return
	}

	conn = &providerConnection{
		info:     prov,
		cConn:    cConn,
		provider: pb.NewProviderClient(cConn),
		store:    pb.NewStorageClient(cConn),
	}

	return
}
