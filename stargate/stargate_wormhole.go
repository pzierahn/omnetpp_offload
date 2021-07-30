package stargate

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type p2pConnector struct {
	conn      *net.UDPConn
	peer      *net.UDPAddr
	mu        sync.RWMutex
	received  bool
	packages  int
	sendDelay time.Duration
}

// Send hello messages to open the NAT
func (p2p *p2pConnector) sendSeeYou(ctx context.Context) (err error) {

	// Wait for some time to ensure that all message get received properly
	timer := time.NewTicker(p2p.sendDelay)
	defer timer.Stop()

	for inx := 0; inx < p2p.packages; inx++ {
		var byt []byte

		p2p.mu.RLock()
		if !p2p.received {
			byt = []byte{0, byte(inx)}
		} else {
			byt = []byte{1, byte(inx)}
		}
		p2p.mu.RUnlock()

		log.Printf("send: %v", byt)

		_, err = p2p.conn.WriteToUDP(byt, p2p.peer)
		if err != nil {
			return
		}

		select {
		case <-timer.C:
			continue
		case <-ctx.Done():
			return
		}
	}

	return
}

// Receive hello messages from peers
func (p2p *p2pConnector) receive() (err error) {

	var br int
	var remote *net.UDPAddr
	var buf = make([]byte, 2)

	var success bool

	for inx := 0; inx < p2p.packages; inx++ {
		br, remote, err = p2p.conn.ReadFromUDP(buf)
		if err != nil {
			return
		}

		if br != 2 {
			log.Printf("received: faulty message...")
			continue
		}

		if remote.String() != p2p.peer.String() {
			log.Printf("received: wrong remote host '%v' should be '%v'",
				remote, p2p.peer)
			continue
		}

		log.Printf("received: peer=%v received=%v\n", p2p.peer, buf[:br])

		p2p.mu.Lock()
		p2p.received = true
		p2p.mu.Unlock()

		if int(buf[1]) == p2p.packages-1 {
			success = buf[0] == 1
			break
		}
	}

	log.Printf("received: success=%v\n", success)

	if !success {
		err = fmt.Errorf("didn't recieve recieve-acknowledgement")
	}

	return
}
