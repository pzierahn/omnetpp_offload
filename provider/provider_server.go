package provider

import (
	"context"
	"github.com/lucas-clemente/quic-go"
	pnet "github.com/patrickz98/project.go.omnetpp/adapter"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/stargate"
	"github.com/patrickz98/project.go.omnetpp/utils"
	"google.golang.org/grpc"
	"log"
)

func Start(conf gconfig.Config) {

	prov := provider{
		providerId: simple.NamedId(conf.Worker.Name, 8),
	}

	log.Printf("start provider (%v)", prov.providerId)

	//
	// Register provider
	//

	qconn, dialer := utils.GRPCDialerAuto()
	log.Printf("quic listener on %v", qconn.LocalAddr())

	brokerConn, err := grpc.Dial(
		conf.Broker.DialAddr(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		//grpc.WithContextDialer(utils.GRPCDialer(conn)),
		grpc.WithContextDialer(dialer),
	)
	if err != nil {
		log.Fatalln(err)
	}

	broker := pb.NewBrokerClient(brokerConn)
	_, err = broker.Register(context.Background(), prov.info())
	if err != nil {
		log.Fatalln(err)
	}

	//
	// Start provider
	//

	log.Println("wait for stargate connection")

	conn, _ := stargate.Connect(prov.providerId)
	defer func() { _ = conn.Close() }()

	tlsConf, _ := utils.GenerateTLSConfig()
	ql, err := quic.Listen(conn, tlsConf, nil)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("create adapter listener")
	lis := pnet.Listen(ql)
	defer func() { _ = lis.Close() }()

	log.Println("listening for consumer")

	server := grpc.NewServer()
	pb.RegisterProviderServer(server, &prov)
	err = server.Serve(lis)

	return
}
