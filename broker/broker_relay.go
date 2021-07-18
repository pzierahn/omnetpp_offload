package broker

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"log"
	"sync"
)

type relayService struct {
	pb.UnimplementedStargateServer
	mu    sync.RWMutex
	match map[string]chan *pb.Port
}

func (relay *relayService) Relay(ctx context.Context, req *pb.RelayRequest) (addr *pb.Port, err error) {

	log.Printf("Relay: DialAddr=%v", req.DialAddr)

	relay.mu.RLock()
	ch, ok := relay.match[req.DialAddr]
	relay.mu.RUnlock()

	if ok {

		//
		// Both clients try to connect
		//

		var port1 int
		var port2 int

		port1, port2, err = stargate.RelayServerTCP()
		if err != nil {
			return
		}

		ch <- &pb.Port{
			Port: uint32(port1),
		}

		addr = &pb.Port{
			Port: uint32(port2),
		}
	} else {

		//
		// Wait for other client
		//

		ch = make(chan *pb.Port)
		defer close(ch)

		relay.mu.Lock()
		relay.match[req.DialAddr] = ch
		relay.mu.Unlock()

		defer func() {
			relay.mu.Lock()
			delete(relay.match, req.DialAddr)
			relay.mu.Unlock()
		}()

		select {
		case addr = <-ch:
		case <-ctx.Done():
			addr = &pb.Port{}
		}
	}

	return
}
