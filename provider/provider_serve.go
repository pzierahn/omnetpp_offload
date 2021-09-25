package provider

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/mimic"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
	"net"
)

func (prov *provider) serve(listener net.Listener) {
	server := grpc.NewServer()
	pb.RegisterProviderServer(server, prov)
	pb.RegisterStorageServer(server, prov.store)
	err := server.Serve(listener)
	if err != nil {
		log.Println(err)
	}
}

func (prov *provider) serveRelay(ctx context.Context) {
	for {
		conn, err := stargate.DialRelayTCP(ctx, prov.providerId)
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			log.Printf("serveRelay: listening LocalAddr=%v RemoteAddr=%v",
				conn.LocalAddr(), conn.RemoteAddr())

			defer func() {
				_ = conn.Close()
			}()

			listener := mimic.TCPConnToListener(conn)
			prov.serve(listener)
		}()
	}
}

func (prov *provider) serveP2P(ctx context.Context) {
	for {
		p2p, _, err := stargate.DialP2PUDP(ctx, prov.providerId)
		if err != nil {
			log.Println(err)
			continue
		}

		lis, err := mimic.NewQUICListener(p2p)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("serveP2P: new connection addr=%v", lis.Addr())

		go func(lis net.Listener) {

			defer func() {
				_ = lis.Close()
				_ = p2p.Close()
			}()

			prov.serve(lis)
		}(lis)
	}
}

func (prov *provider) serveLocal(ctx context.Context) {

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{})
	if err != nil {
		log.Fatalln(err)
	}

	addr, ok := lis.Addr().(*net.TCPAddr)
	if !ok {
		log.Fatalf("serveLocal: could not cast lis.Addr().(*net.tcpAddr)")
	}

	log.Printf("serveLocal: addr=%v", addr)

	go func() {
		err = stargate.BroadcastTCP(ctx, prov.providerId, addr)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	prov.serve(lis)
}
