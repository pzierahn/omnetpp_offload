package stargate

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func dialAddr(ctx context.Context, conn *net.UDPConn, connectionId string) (remote *net.UDPAddr, err error) {
	log.Printf("dialAddr: connectionId=%s rendezvousAddr=%v conn=%v",
		connectionId, rendezvousAddr, conn.LocalAddr())

	receiveConnect := make(chan bool)
	defer close(receiveConnect)

	go func() {
		buf := make([]byte, 1024)
		var read int

		for {
			read, err = conn.Read(buf)
			if err != nil {
				log.Fatalln(err)
			}

			msg := string(buf[0:read])
			if msg == "hello" {
				continue
			}

			remote, err = net.ResolveUDPAddr("udp", msg)
			if err != nil {
				log.Fatalln(err)
			}

			break
		}

		receiveConnect <- true
	}()

	// send initial udp package to stun server
	log.Printf("send stun signal (%s)", connectionId)
	_, err = conn.WriteTo([]byte(connectionId), rendezvousAddr)
	if err != nil {
		return
	}

	// Send ever 30 seconds a stun signal to the broker keep the NAT open
	signalTick := time.NewTicker(time.Second * 20)
	defer signalTick.Stop()

	go func() {
		for range signalTick.C {
			log.Printf("send stun signal (%s)", connectionId)
			_, err = conn.WriteTo([]byte(connectionId), rendezvousAddr)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}()

	select {
	case <-receiveConnect:
	case <-ctx.Done():
		_ = conn.Close()
		err = fmt.Errorf("error: rendezvou server did not responde in time")
	}

	return
}

func Dial(ctx context.Context, connectionId string) (conn *net.UDPConn, remote *net.UDPAddr, err error) {

	conn, err = net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}

	// Get counterpart udp address
	remote, err = dialAddr(ctx, conn, connectionId)
	if err != nil {
		return
	}

	log.Printf("connecting to %v", remote)

	sendPings := 2

	var wg sync.WaitGroup
	wg.Add(sendPings)

	go func() {
		for inx := 0; inx < sendPings; inx++ {
			message := fmt.Sprintf("hello %d", inx)
			w, err := conn.WriteToUDP([]byte(message), remote)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("send message='%s' (%d bytes)\n", message, w)

			// Wait for to seconds to ensure nat hole punching works!
			if inx == 0 {
				time.Sleep(time.Second * 2)
			}

			wg.Done()
		}
	}()

	done := make(chan bool)
	defer close(done)

	go func() {
		for {
			buf := make([]byte, 1024)
			read, remote, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Println(err)
				return
			}

			// TODO: Check for corrupt messages
			msg := string(buf[0:read])

			if !strings.HasPrefix("hello", msg) {
				//
				// faulty message received
				//
				continue
			}

			log.Printf("received message='%s' from %v\n", msg, remote)

			done <- true
			break
		}
	}()

	select {
	case <-done:
		// everything is okay
		log.Printf("connected to %v", remote)

	case <-ctx.Done():
		// connection timeout: close stuff
		_ = conn.Close()
		err = fmt.Errorf("error: could not establish connection to %v in time", remote)
		return
	}

	wg.Wait()

	return
}
