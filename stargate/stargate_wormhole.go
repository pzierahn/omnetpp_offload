package stargate

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type dialer struct {
	conn *net.UDPConn
	peer *net.UDPAddr
}

// Send hello messages to open the NAT
func (dialer *dialer) sendHellos() (err error) {

	for inx := 0; inx < 2; inx++ {
		message := fmt.Sprintf("hello %d", inx)
		_, err = dialer.conn.WriteToUDP([]byte(message), dialer.peer)
		if err != nil {
			return
		}

		log.Printf("send: message='%s' peer=%v\n", message, dialer.peer)

		if inx == 0 {
			// Wait for two seconds to ensure that all message get received properly
			time.Sleep(time.Second * 2)
		}
	}

	return
}

// Receive hello messages from peers
func (dialer *dialer) receive() (err error) {

	var br int
	var remote *net.UDPAddr
	var buf = make([]byte, 1024)

	for {
		br, remote, err = dialer.conn.ReadFromUDP(buf)
		if err != nil {
			return
		}

		msg := string(buf[0:br])

		// Check for corrupt messages
		if !strings.HasPrefix(msg, "hello") {
			continue
		}

		if remote.String() != dialer.peer.String() {
			log.Printf("received: wrong remote host '%v' should be '%v' --> msg='%v'\n",
				remote, dialer.peer, msg)
			continue
		}

		log.Printf("received: peer=%v message='%s'\n", dialer.peer, msg)

		// Wait for second message before leaving read loop
		if msg == "hello 1" {
			return
		}
	}
}
