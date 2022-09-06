package stargrpc

import (
	"context"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/stargate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"
)

// DialLocal creates a gRPC client connection over the local network.
func DialLocal(ctx context.Context, addr stargate.DialAddr) (cc *grpc.ClientConn, err error) {

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
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
}

// ServeLocal takes as an input a stargate dial address and a gRPC server.
func ServeLocal(addr stargate.DialAddr, server *grpc.Server) {

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
