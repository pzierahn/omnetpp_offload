package stargate

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"time"
)

type peerResolver struct {
	conn *net.UDPConn
	dial DialAddr
}

func (client *peerResolver) receivePeer() (peer PeerResolve, err error) {
	buf := make([]byte, 1024)
	var read int

	for {
		read, err = client.conn.Read(buf)
		if err != nil {
			return
		}

		msg := string(buf[0:read])

		log.Printf("receivePeer: msg=%v", msg)

		if msg == "heartbeat" {
			continue
		}

		err = json.Unmarshal(buf[0:read], &peer)

		return
	}
}

func (client *peerResolver) sendDialAddr(ctx context.Context) (err error) {

	// Send a heartbeat signal ever 20 seconds to the broker keep the NAT gate open
	tick := time.NewTicker(time.Second * 20)

	for {
		log.Printf("sendDialAddr: registration signal (dial=%s) to %v", client.dial, rendezvousAddr)

		_, err = client.conn.WriteTo([]byte(client.dial), rendezvousAddr)
		if err != nil {
			log.Printf("################# error: %v", err)
			return
		}

		select {
		case <-tick.C:
		case <-ctx.Done():
			return
		}
	}
}

func (client *peerResolver) resolvePeer(ctx context.Context) (peer PeerResolve, err error) {
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
		peer, recRrr = client.receivePeer()
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
		sendErr = client.sendDialAddr(sendCtx)
		if sendErr != nil {
			once.Do(func() {
				err = sendErr
			})
		}
	}()

	wg.Wait()

	return
}
