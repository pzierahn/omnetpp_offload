package stargate

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func receiveRemoteAddr(ctx context.Context, conn *net.UDPConn, connectionId string) (remote *net.UDPAddr, err error) {
	log.Printf("receiveRemoteAddr: connectionId=%s rendezvousAddr=%v conn=%v",
		connectionId, rendezvousAddr, conn.LocalAddr())

	ch := make(chan error)
	defer close(ch)

	go func() {
		buf := make([]byte, 1024)
		var read int

		for {
			read, err = conn.Read(buf)
			if err != nil {
				ch <- err
				return
			}

			msg := string(buf[0:read])
			if msg == "heartbeat" {
				continue
			}

			remote, err = net.ResolveUDPAddr("udp", msg)
			if err != nil {
				ch <- err
				return
			}

			break
		}

		ch <- nil
	}()

	// Send a heartbeat signal ever 20 seconds to the broker keep the NAT gate open
	signalTick := time.NewTicker(time.Second * 20)
	defer signalTick.Stop()

	go func() {
		for {
			log.Printf("send stun signal connectionId=%s", connectionId)
			_, err = conn.WriteTo([]byte(connectionId), rendezvousAddr)
			if err != nil {
				log.Println(err)
			}

			if _, ok := <-signalTick.C; !ok {
				break
			}
		}
	}()

	select {
	case rec := <-ch:
		if err == nil {
			err = rec
		}
	case <-ctx.Done():
		_ = conn.Close()
		err = fmt.Errorf("error: rendezvou server did not responde in time")
	}

	return
}

func DialUDP(ctx context.Context, connectionId string) (conn *net.UDPConn, remote *net.UDPAddr, err error) {

	conn, err = net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}

	// Get counterparts udp address
	remote, err = receiveRemoteAddr(ctx, conn, connectionId)
	if err != nil {
		return
	}

	log.Printf("connecting to %v", remote)

	hellos := 2

	errCh := make(chan error)
	defer close(errCh)

	go func() {

		//
		// Send hello messages to open the NAT
		//

		for inx := 0; inx < hellos; inx++ {
			message := fmt.Sprintf("hello %d", inx)
			_, err = conn.WriteToUDP([]byte(message), remote)
			if err != nil {
				errCh <- err
				return
			}

			log.Printf("send: message='%s' remote=%v\n", message, remote)

			// Wait for to seconds to ensure NAT hole punching works!
			if inx == 0 {
				time.Sleep(time.Second * 2)
			}
		}

		errCh <- nil
	}()

	go func() {

		//
		// Receive hello messages from peer
		//

		for {
			buf := make([]byte, 1024)
			br, remote, err := conn.ReadFromUDP(buf)
			if err != nil {
				errCh <- err
				return
			}

			msg := string(buf[0:br])

			// Check for corrupt messages
			if !strings.HasPrefix(msg, "hello") {
				continue
			}

			log.Printf("received: message='%s' from %v\n", msg, remote)
			break
		}

		errCh <- nil
	}()

	debugId := rand.Uint32()

	for inx := 0; inx < 2; inx++ {
		select {
		//
		// Receive error or success messages from sender and receiver
		//
		case rec := <-errCh:
			log.Printf("%x error: '%v' receivedErr='%v'", debugId, err, rec)

			if err == nil && rec != nil {
				err = rec
				_ = conn.Close()
			}

		//
		// Connection timeout: close stuff
		//
		case <-ctx.Done():
			_ = conn.Close()
			err = fmt.Errorf("error: could not establish connection to %v in time", remote)
			return
		}
	}

	if err != nil {
		return
	}

	// Everything is okay
	log.Printf("connected to %v", remote)

	return
}
