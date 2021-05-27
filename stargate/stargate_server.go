package stargate

import (
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	pnet "github.com/patrickz98/project.go.omnetpp/net"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type clientType map[string]*net.UDPAddr

var clients = clientType{}

func (c clientType) keys(filter string) (addrs []*net.UDPAddr) {

	for key, val := range c {
		if key != filter {
			addrs = append(addrs, val)
		}
	}

	return
}

func Server() {
	localAddress := ":9595"
	if len(os.Args) > 2 {
		localAddress = os.Args[2]
	}

	addr, _ := net.ResolveUDPAddr("udp", localAddress)
	conn, _ := net.ListenUDP("udp", addr)

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	ql, err := quic.Listen(conn, tlsConf, nil)
	if err != nil {
		log.Fatalln(err)
	}

	lis := pnet.Listen(ql)

	server := grpc.NewServer()
	//pb.RegisterBrokerServer(server, &brk)
	pb.RegisterStorageServer(server, &storage.Server{})
	err = server.Serve(lis)

	//fmt.Printf("addr=%v localAddress=%v\n", addr, localAddress)
	//
	//qconn := quick.Connection{
	//	Connection: conn,
	//}
	//
	//qconn.Init()
	//
	////for {
	//
	//var incoming string
	//_, err := qconn.Receive(&incoming)
	//if err != nil {
	//	panic(err)
	//}
	//
	////log.Printf("incoming '%v'", incoming)
	//log.Printf("incoming: %v bytes", len(incoming))

	//buffer := make([]byte, 1024)
	//bytesRead, remoteAddr, err := conn.ReadFromUDP(buffer)
	//if err != nil {
	//	panic(err)
	//}
	//
	//incoming := string(buffer[0:bytesRead])
	//fmt.Println("[INCOMING]", incoming, remoteAddr)
	//if incoming != "register" {
	//	continue
	//}

	//_, err = conn.WriteToUDP([]byte("Ping from sever"), remoteAddr)
	//if err != nil {
	//	panic(err)
	//}

	//clients[strings.TrimSpace(remoteAddr.String())] = remoteAddr
	//
	//for client := range clients {
	//	resp := clients.keys(client)
	//	if len(resp) > 0 {
	//		var remote *net.UDPAddr
	//		remote, err = net.ResolveUDPAddr("udp", client)
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//
	//		quickConn := quick.Connection{
	//			Connection: conn,
	//		}
	//
	//		garbage := simple.RandomId(1024 * 1024 * 2)
	//		log.Printf("resp: size=%v", len(garbage))
	//		err = quickConn.Send(garbage, remote)
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//
	//		//log.Printf("resp: %v", resp)
	//		//err = quickConn.Send(resp, remote)
	//		//if err != nil {
	//		//	log.Fatalln(err)
	//		//}
	//
	//		log.Printf("[INFO] Responded to %s with %s", client, resp)
	//	}
	//}
	//}
}

// GOOS=linux GOARCH=amd64 go build cmd/main.go
// scp main 7zierahn@rzssh1.informatik.uni-hamburg.de:~
