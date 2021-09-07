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
	packages  int           // Set how many receive acknowledgments should be send
	sendDelay time.Duration // Set a delay for sending receive acknowledgments
	timeout   time.Duration // Set a timeout for establishing the connection and exchanging messages
}

//type receiveAck []byte
//
//func (ack receiveAck) received() bool {
//	return ack[0] == 1
//}

// Send hello messages to open the NAT
func (p2p *p2pConnector) sendSeeYou(ctx context.Context) (err error) {

	// Wait for some time to ensure that all message get received properly
	timer := time.NewTicker(p2p.sendDelay)
	defer timer.Stop()

	for inx := 0; inx < p2p.packages; inx++ {
		var ack []byte

		p2p.mu.RLock()
		if !p2p.received {
			ack = []byte{0, byte(inx)}
		} else {
			ack = []byte{1, byte(inx)}
		}
		p2p.mu.RUnlock()

		log.Printf("send: %v", ack)

		_, err = p2p.conn.WriteToUDP(ack, p2p.peer)
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
func (p2p *p2pConnector) receive() error {

	var ack = make([]byte, 2)
	var success bool

	for inx := 0; inx < p2p.packages; inx++ {
		br, remote, err := p2p.conn.ReadFromUDP(ack)
		if err != nil {
			return err
		}

		if br != 2 {
			continue
		}

		if remote.String() != p2p.peer.String() {
			continue
		}

		log.Printf("received: peer=%v received=%v\n", p2p.peer, ack[:br])

		p2p.mu.Lock()
		p2p.received = true
		p2p.mu.Unlock()

		if int(ack[1]) == p2p.packages-1 {
			success = ack[0] == 1
			break
		}
	}

	if !success {
		return fmt.Errorf("p2p connect failed: no recieve-acknowledgement")
	}

	return nil
}

func (p2p *p2pConnector) connect(ctx context.Context) (err error) {

	ctx, cnl := context.WithTimeout(ctx, p2p.timeout)
	defer cnl()

	if deadline, ok := ctx.Deadline(); ok {
		err = p2p.conn.SetReadDeadline(deadline)
		if err != nil {
			return
		}

		defer func() {
			// Reset deadline
			_ = p2p.conn.SetDeadline(time.Time{})
		}()
	}

	var wg sync.WaitGroup
	var once sync.Once

	wg.Add(1)
	go func() {
		defer wg.Done()

		receiveErr := p2p.receive()
		if receiveErr != nil {
			once.Do(func() {
				err = receiveErr
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		sendErr := p2p.sendSeeYou(ctx)
		if sendErr != nil {
			once.Do(func() {
				err = sendErr
			})
		}
	}()

	wg.Wait()

	return
}
