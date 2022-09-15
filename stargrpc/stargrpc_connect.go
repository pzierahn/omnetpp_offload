package stargrpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

const (
	ConnectLocal = 1 << iota
	ConnectP2P
	ConnectRelay
	ConnectAll = ConnectLocal | ConnectP2P | ConnectRelay
)

func ConnectionToName(connection int) (name string) {
	if connection&ConnectLocal != 0 {
		return "local"
	}

	if connection&ConnectP2P != 0 {
		return "p2p"
	}

	if connection&ConnectRelay != 0 {
		return "relay"
	}

	return "none"
}

func NameToConnection(name string) (connection int) {
	switch name {
	case "local":
		return ConnectLocal
	case "p2p":
		return ConnectP2P
	case "relay":
		return ConnectRelay
	default:
		return ConnectAll
	}
}

func ConnectFeedback(ctx context.Context, addr string, connect int) (client *grpc.ClientConn, connection int, err error) {

	if connect == 0 {
		connect = ConnectAll
	}

	if connect&ConnectLocal != 0 {
		client, err = DialLocal(ctx, addr)
		if err == nil {
			log.Printf("[%v] Connect over local network", addr)
			return client, ConnectLocal, nil
		}
	}

	if connect&ConnectP2P != 0 {
		client, err = DialP2P(ctx, addr)
		if err == nil {
			log.Printf("[%v] Connect over peer-to-peer", addr)
			return client, ConnectP2P, nil
		}
	}

	if connect&ConnectRelay != 0 {
		client, err = DialRelay(ctx, addr)
		if err == nil {
			log.Printf("[%v] Connect over relay server", addr)
			return client, ConnectRelay, nil
		}
	}

	if err != nil {
		err = fmt.Errorf("couldn't establish a connection: %v", err)
	} else {
		err = fmt.Errorf("couldn't establish a connection")
	}

	return
}

func Connect(ctx context.Context, addr string, connect int) (client *grpc.ClientConn, err error) {
	client, _, err = ConnectFeedback(ctx, addr, connect)
	return
}
