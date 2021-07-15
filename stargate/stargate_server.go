package stargate

import (
	"context"
	"net"
	"sync"
	"time"
)

//type dialAddr = string

type stargateServer struct {
	ctx        context.Context
	conn       *net.UDPConn
	mu         sync.RWMutex
	rendezvous map[string]map[string]*net.UDPAddr
	timers     map[string]*time.Timer
}

func (server *stargateServer) DebugValues() (values map[string][]string) {

	server.mu.RLock()
	defer server.mu.RUnlock()

	values = make(map[string][]string)
	for id, val := range server.rendezvous {
		for conn := range val {
			values[id] = append(values[id], conn)
		}
	}

	return
}

func (server *stargateServer) heartbeatDispatcher() {

	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			server.mu.RLock()

			for id, register := range server.rendezvous {
				for _, addr := range register {
					log.Printf("send heartbeat: connectionId=%v addr=%v", id, addr)

					_, err := server.conn.WriteTo([]byte("hello"), addr)
					if err != nil {
						log.Println(err)
					}
				}
			}

			server.mu.RUnlock()

		case <-server.ctx.Done():
			//
			// Exit
			//

			break
		}
	}
}

func (server *stargateServer) resetCleaner(dialAddr, remoteAddr string) {
	timerKey := dialAddr + "-" + remoteAddr

	wait := time.Second * 40
	timer, ok := server.timers[timerKey]

	if ok {

		//
		// Reset existing timer
		//

		timer.Reset(wait)
		return
	}

	//
	// Create a new timer to remove obsolete data
	//

	timer = time.NewTimer(wait)
	server.timers[timerKey] = timer

	// Create a new go-routine that will delete the data
	go func() {
		if _, ok := <-timer.C; ok {

			server.mu.Lock()
			defer server.mu.Unlock()

			log.Printf("remove old connection trace connectId=%s", dialAddr)

			delete(server.timers, timerKey)

			delete(server.rendezvous[dialAddr], remoteAddr)

			if len(server.rendezvous[dialAddr]) == 0 {
				delete(server.rendezvous, dialAddr)
			}
		}
	}()
}

func (server *stargateServer) receiveDial() (err error) {

	buffer := make([]byte, 1024)

	br, remoteAddr, err := server.conn.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	dialAddr := string(buffer[0:br])
	log.Printf("receive: dialAddr=%s remoteAddr=%v", dialAddr, remoteAddr)

	server.mu.Lock()
	defer server.mu.Unlock()

	rendezvous := server.rendezvous

	if _, ok := rendezvous[dialAddr]; !ok {
		rendezvous[dialAddr] = make(map[string]*net.UDPAddr)
	}

	server.resetCleaner(dialAddr, remoteAddr.String())

	rendezvous[dialAddr][remoteAddr.String()] = remoteAddr

	if len(rendezvous[dialAddr]) != 2 {
		return
	}

	hosts := make([]*net.UDPAddr, 2)

	var inx int
	for _, host := range rendezvous[dialAddr] {
		hosts[inx] = host
		inx++
	}

	_, err = server.conn.WriteToUDP([]byte(hosts[0].String()), hosts[1])
	if err != nil {
		log.Fatalln(err)
	}

	_, err = server.conn.WriteToUDP([]byte(hosts[1].String()), hosts[0])
	if err != nil {
		log.Fatalln(err)
	}

	delete(rendezvous, dialAddr)

	return
}

func Server(ctx context.Context) (err error) {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: stunPort})
	if err != nil {
		return
	}

	log.Printf("start stargate server on %v", conn.LocalAddr())

	server := stargateServer{
		ctx:        ctx,
		conn:       conn,
		rendezvous: make(map[string]map[string]*net.UDPAddr),
		timers:     make(map[string]*time.Timer),
	}

	go server.heartbeatDispatcher()

	for {
		err = server.receiveDial()
		if err != nil {
			log.Println(err)
		}
	}
}

// GOOS=linux GOARCH=amd64 go build cmd/main.go
// scp main 7zierahn@rzssh1.informatik.uni-hamburg.de:~
