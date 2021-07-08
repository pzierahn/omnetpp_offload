package consumer

import (
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/equic"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type connection struct {
	simulation *pb.Simulation
	info       *pb.ProviderInfo
	conn       *net.UDPConn
	cConn      *grpc.ClientConn
	provider   pb.ProviderClient
	store      pb.StorageClient
}

func (conn *connection) name() (name string) {
	return fmt.Sprintf("%-20s", conn.info.ProviderId)
}

func connect(prov *pb.ProviderInfo) (conn *connection, err error) {

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

	conn = &connection{
		info:     prov,
		conn:     gate,
		cConn:    cConn,
		provider: pb.NewProviderClient(cConn),
		store:    pb.NewStorageClient(cConn),
	}

	return
}

func (conn *connection) close() {
	_ = conn.cConn.Close()
	_ = conn.conn.Close()
}
