package stargate

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	port = 9595
)

var rendezvousAddr = &net.UDPAddr{
	IP:   net.ParseIP("31.18.129.212"),
	Port: port,
}

func Connect(connectionId string) (conn *net.UDPConn, remote *net.UDPAddr) {

	var err error
	conn, err = net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("connectionId=%s rendezvousAddr=%v conn=%v",
		connectionId, rendezvousAddr, conn.LocalAddr())

	bytesWritten, err := conn.WriteTo([]byte(connectionId), rendezvousAddr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("send register connectionId=%s (%d bytes)", connectionId, bytesWritten)

	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		log.Fatalln(err)
	}

	remote, err = net.ResolveUDPAddr("udp", string(buffer[0:bytesRead]))
	if err != nil {
		log.Fatalln(err)
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

	listen(conn)

	wg.Wait()

	return
}

func listen(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	bytesRead, remote, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("receive '%s' from %v\n", string(buffer[0:bytesRead]), remote)
}
