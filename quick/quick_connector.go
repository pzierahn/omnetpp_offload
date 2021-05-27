package quick

import (
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	//payloadSize = 8192
	payloadSize = 1
)

type ping struct {
	sessionId  uint32
	remoteAddr *net.UDPAddr
}

type Connection struct {
	sync.RWMutex
	Connection *net.UDPConn
	connMu     sync.Mutex
	pings      chan<- ping
}

func (conn *Connection) digestPing(pkg parcel, remote *net.UDPAddr) {
	log.Printf("digestPing: SessionId=%x MessageId=%x ", pkg.sessionId, pkg.messageId)

	conn.RLock()
	conn.pings <- ping{
		sessionId:  0,
		remoteAddr: remote,
	}
	conn.RUnlock()
}

func (conn *Connection) Init() {
	conn.pings = make(map[uint32]chan<- ping)

	go func() {
		// max UDP package size (2 ^ 16 = 65536)
		buffer := make([]byte, 65536)

		for {
			bytesRead, remote, err := conn.Connection.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			var pkg parcel
			err = decodeGob(buffer[0:bytesRead], &pkg)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			if pkg.cast == castPing {
				go conn.digestPing(pkg, remote)
			}
		}
	}()
}

func (conn *Connection) Ping(remote *net.UDPAddr) {

	pings := make(chan ping)
	conn.Lock()
	conn.pings = pings
	conn.Unlock()

	//
	// Wait for session
	//

	var lSessionId uint32
	var rSessionId uint32

	sessionLookup := make(map[uint32]uint32)

	lSessionId = rand.Uint32()

	pingTic := time.NewTicker(time.Millisecond * 50)

	go func() {

	}()

	for range pingTic.C {
		log.Printf("send ping")
	}
}
