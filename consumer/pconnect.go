package consumer

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/eval"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargrpc"
	"google.golang.org/grpc"
	"log"
	"os"
)

const (
	connectLocal = 1 << iota
	connectP2P
	connectRelay
	connectAll = connectLocal | connectP2P | connectRelay
)

type providerConnection struct {
	ctx           context.Context
	conn          *grpc.ClientConn
	info          *pb.ProviderInfo
	provider      pb.ProviderClient
	store         pb.StorageClient
	downloadQueue chan *download
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

	close(pConn.downloadQueue)
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
		conn, err = stargrpc.DialLocal(ctx, prov.ProviderId)
		if err == nil {
			eval.LogSetup(eval.ConnectLocal, prov)
			return
		}
	}

	if connect&connectP2P != 0 {
		conn, err = stargrpc.DialP2P(ctx, prov.ProviderId)
		if err == nil {
			eval.LogSetup(eval.ConnectP2P, prov)
			return
		}
	}

	if connect&connectRelay != 0 {
		conn, err = stargrpc.DialRelay(ctx, prov.ProviderId)
		if err == nil {
			eval.LogSetup(eval.ConnectRelay, prov)
			return
		}
	}

	return
}
