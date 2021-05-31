package broker

import (
	"github.com/lucas-clemente/quic-go"
	pnet "github.com/pzierahn/project.go.omnetpp/adapter"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"github.com/pzierahn/project.go.omnetpp/utils"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Start(conf Config) (err error) {

	go stargate.Server()

	log.Printf("start server on :%d", conf.BrokerPort)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: conf.BrokerPort,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = conn.Close() }()

	tlsConf, _ := utils.GenerateTLSConfig()

	ql, err := quic.Listen(conn, tlsConf, nil)
	if err != nil {
		log.Fatalln(err)
	}

	lis := pnet.Listen(ql)
	defer func() { _ = lis.Close() }()

	brk := broker{
		providers:   make(map[string]*pb.ProviderInfo),
		utilization: make(map[string]*pb.Utilization),
	}

	go brk.startWebService()

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &brk)
	pb.RegisterStorageServer(server, &storage.Server{})
	err = server.Serve(lis)

	return
}
