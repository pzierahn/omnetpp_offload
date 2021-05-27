package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/patrickz98/project.go.omnetpp/adapter"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"google.golang.org/grpc"
	"net"
)

type workerConnection struct {
	providerId string
	config     gconfig.Worker
	conn       *grpc.ClientConn
	broker     pb.BrokerClient
	storage    storage.Client
	agents     int
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

	dialer := grpc.WithContextDialer(func(ctx context.Context, target string) (conn net.Conn, err error) {

		logger.Printf("############# target=%v", target)

		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"quic-echo-example"},
		}
		sess, err := quic.DialAddrContext(ctx, target, tlsConf, &quic.Config{
			KeepAlive: true,
		})
		if err != nil {
			logger.Printf("############# err: %v", err)
			return
		}

		stream, err := sess.OpenStreamSync(ctx)
		if err != nil {
			logger.Printf("############# err: %v", err)
			return
		}

		conn = &adapter.Conn{Sess: sess, Stream: stream}
		logger.Printf("############# conn=%v", conn)

		return

	})

	conn, err := grpc.Dial(config.Broker.DialAddr(), grpc.WithInsecure(), grpc.WithBlock(), dialer)
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		return
	}

	logger.Printf("############# connected!")

	worker = &workerConnection{
		providerId: simple.NamedId(config.Worker.Name, 8),
		conn:       conn,
		broker:     pb.NewBrokerClient(conn),
		//storage:    storage.InitClient(config.Broker),
		agents:     config.Worker.DevoteCPUs,
	}

	return
}
