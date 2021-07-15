package stargate

import (
	"context"
	"net"
	"sync"
	"time"
)

//type dialAddr = string

const (
	cleanTimeout = time.Second * 40
)

type waiter struct {
	addr    *net.UDPAddr
	timeout *time.Timer
}

type stargateServer struct {
	conn       *net.UDPConn
	mu         sync.RWMutex
	waiting    map[string]*waiter
	rendezvous map[string]chan *net.UDPAddr
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

			break
		}
	}
}

func (server *stargateServer) prune(dialAddr string, addr *net.UDPAddr) {
	server.mu.Lock()
	defer server.mu.Unlock()

	log.Printf("pruning: dialAddr=%v addr=%v", dialAddr, addr)

	if ch, ok := server.rendezvous[dialAddr]; ok {
		delete(server.rendezvous, dialAddr)
		close(ch)
	}

	delete(server.waiting, addr.String())
}

func (server *stargateServer) receiveDial() (err error) {

	buffer := make([]byte, 1024)

	br, addr, err := server.conn.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	dialAddr := string(buffer[0:br])
	log.Printf("receive: dialAddr=%s remoteAddr=%v", dialAddr, addr)

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

	if ch, ok := server.rendezvous[dialAddr]; ok {
		//
		// Other peer already waiting
		//

		peerAddr := <-ch
		defer func() {
			delete(server.rendezvous, dialAddr)
			close(ch)

			delete(server.waiting, peerAddr.String())
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
		// Waiting for peer to dial in
		//

		ch = make(chan *net.UDPAddr, 1)
		timeout := time.NewTimer(cleanTimeout)

		server.rendezvous[dialAddr] = ch
		server.waiting[addr.String()] = &waiter{
			addr:    addr,
			timeout: timeout,
		}

		go func() {
			if _, prune := <-timeout.C; prune {
				server.prune(dialAddr, addr)
			}
		}()

		ch <- addr
	}

	return
}

func Server(ctx context.Context) (err error) {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: stunPort})
	if err != nil {
		return
	}

	log.Printf("start stargate server on %v", conn.LocalAddr())

	server := stargateServer{
		conn:       conn,
		waiting:    make(map[string]*waiter),
		rendezvous: make(map[string]chan *net.UDPAddr),
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
