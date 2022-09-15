package consumer

import (
	"context"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/stargrpc"
	"google.golang.org/grpc"
	"log"
	"os"
	"sync"
)

type providerConnection struct {
	ctx        context.Context
	client     *grpc.ClientConn
	connection int
	info       *pb.ProviderInfo
	provider   pb.ProviderClient
	store      pb.StorageClient
	dmu        sync.Mutex
}

type download struct {
	task *pb.SimulationRun
	ref  *pb.StorageRef
}

func (connect *providerConnection) id() (name string) {
	return connect.info.ProviderId
}

func (connect *providerConnection) close() {
	//TODO: connect.provider.DropSession(ctx, &pb.Session{})

	//close(connect.downloadQueue)
	_ = connect.client.Close()
}

func pconnect(ctx context.Context, prov *pb.ProviderInfo) (pconn *providerConnection, err error) {

	connect := stargrpc.ConnectAll

	// Eval stuff to ensure that only the desired connection will be used
	switch os.Getenv("CONNECT") {
	case "local":
		log.Println("########################## eval debug: connect only local!")
		connect = stargrpc.ConnectLocal

	case "p2p":
		log.Println("########################## eval debug: connect only p2p!")
		connect = stargrpc.ConnectP2P

	case "relay":
		log.Println("########################## eval debug: connect only relay!")
		connect = stargrpc.ConnectRelay
	}

	client, conn, err := stargrpc.ConnectFeedback(ctx, prov.ProviderId, connect)
	if err != nil {
		return nil, err
	}

	pconn = &providerConnection{
		client:     client,
		connection: conn,
		info:       prov,
		provider:   pb.NewProviderClient(client),
		store:      pb.NewStorageClient(client),
	}

	return
}
