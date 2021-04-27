package worker

import (
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"google.golang.org/grpc"
)

type workerConnection struct {
	config        Config
	conn          *grpc.ClientConn
	client        pb.BrokerClient
	freeResources int
}

func (client *workerConnection) Close() (err error) {
	err = client.conn.Close()
	return
}

func Connect(config Config) (worker *workerConnection, err error) {

	config.workerId = simple.NamedId(config.WorkerName, 8)

	//
	// Setup a connection to the server
	//

	conn, err := grpc.Dial(config.Broker.DialAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		return
	}

	client := pb.NewBrokerClient(conn)

	worker = &workerConnection{
		config:        config,
		conn:          conn,
		client:        client,
		freeResources: config.DevoteCPUs,
	}

	return
}
