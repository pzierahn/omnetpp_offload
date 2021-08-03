package provider

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/equic"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
	"net"
)

func (prov *provider) listenRelay() {

	for {
		ctx := context.Background()
		conn, err := stargate.DialRelayTCP(ctx, prov.providerId)
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			log.Printf("listenRelay: listening LocalAddr=%v RemoteAddr=%v",
				conn.LocalAddr(), conn.RemoteAddr())

			server := grpc.NewServer()
			pb.RegisterProviderServer(server, prov)
			pb.RegisterStorageServer(server, prov.store)
			err := server.Serve(equic.ListenTCP(conn))
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func (prov *provider) listenP2P() {
	for {
		ctx := context.Background()
		p2p, err := equic.P2PListener(ctx, prov.providerId)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("listenP2P: new connection addr=%v", p2p.Addr())

		go func(p2p net.Listener) {

			// TODO: Find a way to close the p2p connection properly
			defer func() { _ = p2p.Close() }()

			server := grpc.NewServer()
			pb.RegisterProviderServer(server, prov)
			pb.RegisterStorageServer(server, prov.store)
			err := server.Serve(p2p)
			if err != nil {
				log.Println(err)
			}
		}(p2p)
	}
}

func (prov *provider) listenLocal() {

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{})
	if err != nil {
		log.Fatalln(err)
	}

	addr, ok := lis.Addr().(*net.TCPAddr)
	if !ok {
		log.Fatalf("listenLocal: could not cast lis.Addr().(*net.tcpAddr)")
	}

	log.Printf("listenLocal: addr=%v", addr)

	go func() {
		ctx := context.Background()
		err = stargate.BroadcastTCP(ctx, prov.providerId, addr)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	server := grpc.NewServer()
	pb.RegisterProviderServer(server, prov)
	pb.RegisterStorageServer(server, prov.store)
	if err != server.Serve(lis) {
		log.Println(err)
	}
}
