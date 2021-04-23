package worker

import (
	"context"
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"runtime"
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
	logger.Println("config", simple.PrettyString(config))

	//
	// Setup a connection to the server
	//

	conn, err := grpc.Dial(config.BrokerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		return
	}

	client := pb.NewBrokerClient(conn)

	md := metadata.New(map[string]string{
		"workerId": config.WorkerId,
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"numCPU":   fmt.Sprint(runtime.NumCPU()),
	})

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	worker = &workerConnection{
		config:        config,
		conn:          conn,
		client:        client,
		freeResources: runtime.NumCPU(),
	}

	return
}
