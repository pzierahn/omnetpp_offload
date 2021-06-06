package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"github.com/pzierahn/project.go.omnetpp/utils"
	"google.golang.org/grpc"
	"log"
	"time"
)

type connection struct {
	info     *pb.ProviderInfo
	provider pb.ProviderClient
	store    pb.StorageClient
}

func (cons *consumer) connect(prov *pb.ProviderInfo) (conn *connection, err error) {

	log.Printf("connect to provider (%v)", prov.ProviderId)

	ctx, cln := context.WithTimeout(context.Background(), time.Second*4)
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
		grpc.WithContextDialer(utils.GRPCDialer(gate)),
	)
	if err != nil {
		return
	}

	conn = &connection{
		info:     prov,
		provider: pb.NewProviderClient(cConn),
		store:    pb.NewStorageClient(cConn),
	}

	cons.connMu.Lock()
	cons.connections[prov.ProviderId] = conn
	cons.connMu.Unlock()

	// TODO: Handle disconnect!

	return
}
