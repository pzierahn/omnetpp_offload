package stargate

import (
	"context"
	"net"
	"sync"
	"time"
)

func DialUDP(ctx context.Context, dialAddr DialAddr) (conn *net.UDPConn, peer *net.UDPAddr, err error) {

	log.Printf("DialUDP: dialAddr=%v", dialAddr)

	conn, err = net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}

	if deadline, ok := ctx.Deadline(); ok {
		log.Printf("DialUDP: deadline=%v", deadline)
		err = conn.SetDeadline(deadline)
		if err != nil {
			return
		}

		defer func() {
			// Reset deadline
			_ = conn.SetDeadline(time.Time{})
		}()
	}

	client := stargateClient{
		conn: conn,
		dial: dialAddr,
	}
	peer, err = client.resolvePeer()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	var once sync.Once

	log.Printf("DialUDP: dialAddr=%v peer=%v", dialAddr, peer)

	helper := dialer{
		conn: conn,
		peer: peer,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		sendErr := helper.sendHellos()
		if sendErr != nil {
			once.Do(func() {
				err = sendErr
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		receiveErr := helper.receive()
		if receiveErr != nil {
			once.Do(func() {
				err = receiveErr
			})
		}
	}()

	wg.Wait()

	if err != nil {
		return
	}

	// Everything is okay
	log.Printf("DialUDP: dialAddr=%v peer=%v connection established", dialAddr, peer)

	return
}
