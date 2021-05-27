package broker

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	pnet "github.com/patrickz98/project.go.omnetpp/adapter"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"google.golang.org/grpc"
	"log"
	"math/big"
	"net"
	"time"
)

type broker struct {
	pb.UnimplementedBrokerServer
	providers   providerManager
	simulations simulationManager
}

func Start(conf Config) (err error) {

	logger.Println("start server on", conf.Port)

	var lis net.Listener
	//lis, err = net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	//if err != nil {
	//	return
	//}

	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", conf.Port))
	conn, _ := net.ListenUDP("udp", addr)

	ql, err := quic.Listen(conn, generateTLSConfig(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	lis = pnet.Listen(ql)

	defer func() { _ = lis.Close() }()

	brk := broker{
		providers:   newProviderManager(),
		simulations: newSimulationManager(),
	}

	if conf.WebInterface {
		go brk.startWebService()
	}

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &brk)
	pb.RegisterStorageServer(server, &storage.Server{})

	go func() {
		for range time.Tick(time.Second * 4) {
			brk.distribute()
		}
	}()

	err = server.Serve(lis)

	return
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
