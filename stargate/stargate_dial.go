package stargate

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func dialAddr(ctx context.Context, conn *net.UDPConn, connectionId string) (remote *net.UDPAddr, err error) {
	log.Printf("connectionId=%s rendezvousAddr=%v conn=%v",
		connectionId, rendezvousAddr, conn.LocalAddr())

	wr, err := conn.WriteTo([]byte(connectionId), rendezvousAddr)
	if err != nil {
		return
	}

	log.Printf("send register connectionId=%s (%d bytes)", connectionId, wr)

	success := make(chan bool)
	defer close(success)

	go func() {
		buffer := make([]byte, 1024)
		var read int

		read, err = conn.Read(buffer)
		if err != nil {
			return
		}

		remote, err = net.ResolveUDPAddr("udp", string(buffer[0:read]))
		success <- true
	}()

	select {
	case <-success:
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

	remote, err = dialAddr(ctx, conn, connectionId)
	if err != nil {
		return
	}

	log.Printf("connect to %v", remote)

	sendPings := 2

	var wg sync.WaitGroup
	wg.Add(sendPings)

	go func() {
		for inx := 0; inx < sendPings; inx++ {
			message := fmt.Sprintf("ping %d", inx)
			w, err := conn.WriteToUDP([]byte(message), remote)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("send '%s' (%d bytes)\n", message, w)

			// Wait for to seconds to ensure nat hole punching works!
			if inx == 0 {
				time.Sleep(time.Second * 2)
			}

			wg.Done()
		}
	}()

	select {
	case <-listen(conn):
		// everything is okay
	case <-ctx.Done():
		// connection timeout: close stuff
		_ = conn.Close()
		err = fmt.Errorf("error: could not establish connection in time")
		return
	}

	wg.Wait()

	return
}

func listen(conn *net.UDPConn) (done chan bool) {

	done = make(chan bool)

	go func() {
		buffer := make([]byte, 1024)
		bytesRead, remote, err := conn.ReadFromUDP(buffer)
		if err != nil {
			//log.Println(err)
			return
		}

		log.Printf("receive '%s' from %v\n", string(buffer[0:bytesRead]), remote)

		done <- true
	}()

	return
}
