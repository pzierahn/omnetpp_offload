package stargate

import (
	"context"
	"fmt"
	"net"
	"time"
)

type p2pConnector struct {
	conn    *net.UDPConn
	start   bool
	peer    *net.UDPAddr
	timeout time.Duration // Set a timeout for establishing the connection and exchanging messages
}

const (
	_ = iota
	msgPenetrateNAT
	msgHello
	msgACK
)

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

	//
	// Send first message to open the NAT. This message is possibly lost.
	//

	log.Printf("Write msgPenetrateNAT")
	_, err = p2p.conn.WriteToUDP([]byte{msgPenetrateNAT}, p2p.peer)
	if err != nil {
		return
	}

	if p2p.start {
		// Sleep to ensure that peer-2 has time to receive the message.
		time.Sleep(time.Millisecond * 30)

		log.Printf("Write msgHello")
		_, err = p2p.conn.WriteToUDP([]byte{msgHello}, p2p.peer)
		if err != nil {
			return
		}
	}

	buf := make([]byte, 1)

	for {
		_, err = p2p.conn.Read(buf)
		if err != nil {
			return
		}

		log.Printf("Read %v", buf)

		if buf[0] == msgPenetrateNAT {

			//
			// Ignore NAT opening message.
			//

			continue
		}

		break
	}

	if buf[0] == msgACK || buf[0] == msgHello {
		log.Printf("Write msgACK")
		_, err = p2p.conn.WriteToUDP([]byte{msgACK}, p2p.peer)
	} else {
		err = fmt.Errorf("received garbage: %v", buf[0])
		return
	}

	if !p2p.start {
		_, err = p2p.conn.Read(buf)
		if err != nil {
			return
		}

		log.Printf("Read %v", buf)

		if buf[0] != msgACK {
			err = fmt.Errorf("peer didn't receive message. %v", buf[0])
		}
	}

	return
}
