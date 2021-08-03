package stargate

import (
	"context"
	"net"
	"sync"
	"time"
)

type peerResolver struct {
	conn *net.UDPConn
	dial DialAddr
}

func (client *peerResolver) receive() (peer *net.UDPAddr, err error) {
	buf := make([]byte, 1024)
	var read int

	for {
		read, err = client.conn.Read(buf)
		if err != nil {
			return
		}

		msg := string(buf[0:read])

		log.Printf("receive: msg='%v'", msg)

		if msg == "heartbeat" {
			continue
		}

		return net.ResolveUDPAddr("udp", msg)
	}
}

func (client *peerResolver) send(ctx context.Context) (err error) {

	// Send a heartbeat signal ever 20 seconds to the broker keep the NAT gate open
	tick := time.NewTicker(time.Second * 20)

	for {
		log.Printf("send: registration signal (dial=%s)", client.dial)

		_, err = client.conn.WriteTo([]byte(client.dial), rendezvousAddr)
		if err != nil {
			return
		}

		select {
		case <-tick.C:
		case <-ctx.Done():
			return
		}
	}
}

func (client *peerResolver) resolvePeer(ctx context.Context) (peer *net.UDPAddr, err error) {
	log.Printf("resolvePeer: dialAddr=%s conn=%v", client.dial, client.conn.LocalAddr())

	if deadline, ok := ctx.Deadline(); ok {

		//
		// Set deadline for peer resolving
		//

		log.Printf("resolvePeer: deadline=%v", deadline)
		err = client.conn.SetDeadline(deadline)
		if err != nil {
			return
		}

		defer func() {
			// Reset deadline
			_ = client.conn.SetDeadline(time.Time{})
		}()
	}

	sendCtx, cnlSend := context.WithCancel(ctx)

	// Wait group to ensure that sending and receiving is finished before returning
	var wg sync.WaitGroup

	// This once will be used to set the initial error
	var once sync.Once

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cnlSend()

		var recRrr error
		peer, recRrr = client.receive()
		if recRrr != nil {
			once.Do(func() {
				err = recRrr
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		var sendErr error
		sendErr = client.send(sendCtx)
		if sendErr != nil {
			once.Do(func() {
				err = sendErr
			})
		}
	}()

	wg.Wait()

	return
}
