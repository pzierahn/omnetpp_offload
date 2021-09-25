package egrpc

import (
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func DialLocal(ctx context.Context, addr string) (cc *grpc.ClientConn, err error) {

	log.Printf("DialLocal: addr=%v", addr)

	ctx, cln := context.WithTimeout(ctx, time.Second)
	defer cln()

	raddr, err := stargate.DialLocal(ctx, addr)
	if err != nil {
		return
	}

	log.Printf("DialLocal: addr=%v dial=%v", addr, raddr)

	return grpc.DialContext(
		ctx,
		raddr.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

func ServeLocal(addr string, server *grpc.Server) {

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{})
	if err != nil {
		return
	}
	defer func() {
		_ = lis.Close()
	}()

	raddr, ok := lis.Addr().(*net.TCPAddr)
	if !ok {
		err = fmt.Errorf("ServeLocal: could not cast lis.Addr().(*net.tcpAddr)")
		return
	}

	log.Printf("ServeLocal: raddr=%v", raddr)

	go func() {
		ctx := context.Background()
		err = stargate.BroadcastTCP(ctx, addr, raddr)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	if err := server.Serve(lis); err != nil {
		log.Fatalln(err)
	}

	return
}
