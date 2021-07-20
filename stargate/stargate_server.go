package stargate

import (
	"context"
	"encoding/json"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"net"
	"sync"
	"time"
)

type DialAddr = string
type udpAddr = string

const (
	cleanTimeout = time.Second * 40
)

type waiter struct {
	addr     *net.UDPAddr
	timeout  *time.Timer
	dialAddr DialAddr
}

type stargateServer struct {
	conn    *net.UDPConn
	mu      sync.RWMutex
	waiting map[udpAddr]*waiter
	peers   map[DialAddr]*net.UDPAddr
}

var server *stargateServer

func DebugValues() (bytes []byte, err error) {

	server.mu.RLock()
	defer server.mu.RUnlock()

	waiting := make(map[udpAddr]string)
	for addr, wait := range server.waiting {
		waiting[addr] = wait.dialAddr
	}

	data := struct {
		LocalAddr net.Addr
		Waiting   map[udpAddr]string
		Peers     map[DialAddr]*net.UDPAddr
	}{
		LocalAddr: server.conn.LocalAddr(),
		Waiting:   waiting,
		Peers:     server.peers,
	}

	return json.MarshalIndent(data, "", "  ")
}

func (server *stargateServer) heartbeatDispatcher(ctx context.Context) {

	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			server.mu.RLock()

			for _, wait := range server.waiting {
				log.Printf("send heartbeat: addr=%v", wait.addr)

				_, err := server.conn.WriteTo([]byte("heartbeat"), wait.addr)
				if err != nil {
					log.Println(err)
				}
			}

			server.mu.RUnlock()

		case <-ctx.Done():
			//
			// Exit
			//

			return
		}
	}
}

func (server *stargateServer) prune(dial DialAddr, addr *net.UDPAddr) {
	server.mu.Lock()
	defer server.mu.Unlock()

	log.Printf("pruning: dialAddr=%v addr=%v", dial, addr)

	delete(server.peers, dial)
	delete(server.waiting, addr.String())
}

func (server *stargateServer) receiveDial() (err error) {

	buffer := make([]byte, 1024)

	br, addr, err := server.conn.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	dial := string(buffer[0:br])
	log.Printf("receive: dialAddr=%s remoteAddr=%v", dial, addr)

	server.mu.Lock()
	defer server.mu.Unlock()

	if wait, ok := server.waiting[addr.String()]; ok {

		//
		// The dialing clients send periodically new dial signals to ensure that the NAT stays open.
		// When this happens the server reset the waiters timeout to prevent pruning.
		//

		wait.timeout.Reset(cleanTimeout)
		return
	}

	if peerAddr, ok := server.peers[dial]; ok {
		//
		// Other peers already waiting
		//

		defer func() {
			delete(server.peers, dial)

			if wait, ok := server.waiting[peerAddr.String()]; ok {
				wait.timeout.Stop()
				delete(server.waiting, peerAddr.String())
			}
		}()

		_, err = server.conn.WriteToUDP([]byte(addr.String()), peerAddr)
		if err != nil {
			return
		}

		_, err = server.conn.WriteToUDP([]byte(peerAddr.String()), addr)
		if err != nil {
			return
		}
	} else {
		//
		// Waiting for peers to dial in
		//

		timeout := time.NewTimer(cleanTimeout)

		server.peers[dial] = addr
		server.waiting[addr.String()] = &waiter{
			addr:     addr,
			timeout:  timeout,
			dialAddr: dial,
		}

		go func() {
			if _, prune := <-timeout.C; prune {
				server.prune(dial, addr)
			}
		}()
	}

	return
}

func Server(ctx context.Context) (err error) {

	//
	// TODO: start ServerRelayTCP
	//

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: gconfig.StargatePort(),
	})
	if err != nil {
		return
	}

	log.Printf("start stargate server on %v", conn.LocalAddr())

	server = &stargateServer{
		conn:    conn,
		waiting: make(map[string]*waiter),
		peers:   make(map[string]*net.UDPAddr),
	}

	go server.heartbeatDispatcher(ctx)

	for {
		err = server.receiveDial()
		if err != nil {
			log.Println(err)
		}
	}
}

// GOOS=linux GOARCH=amd64 go build cmd/main.go
// scp main 7zierahn@rzssh1.informatik.uni-hamburg.de:~
