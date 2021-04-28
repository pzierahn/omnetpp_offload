package worker

import (
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"google.golang.org/grpc"
)

type workerConnection struct {
	workerId      string
	config        gconfig.Worker
	conn          *grpc.ClientConn
	broker        pb.BrokerClient
	storage       storage.Client
	freeResources int
}

func (client *workerConnection) Close() (err error) {
	err = client.conn.Close()
	return
}

func Init(config gconfig.Config) (worker *workerConnection, err error) {

	logger.Printf("connecting to %s\n", config.Broker.DialAddr())

	//
	// Setup a connection to the server
	//

	conn, err := grpc.Dial(config.Broker.DialAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		return
	}

	worker = &workerConnection{
		workerId:      simple.NamedId(config.Worker.Name, 8),
		conn:          conn,
		broker:        pb.NewBrokerClient(conn),
		storage:       storage.InitClient(config.Broker),
		freeResources: config.Worker.DevoteCPUs,
	}

	return
}
