package stargate

import (
	"net"
	"sync"
	"time"
)

var matchMu sync.RWMutex
var match = make(map[string]map[string]*net.UDPAddr)
var timers = make(map[string]*time.Timer)

func DebugValues() map[string]map[string]*net.UDPAddr {
	matchMu.RLock()
	defer matchMu.RUnlock()

	return match
}

func ping(conn *net.UDPConn) {
	matchMu.RLock()
	defer matchMu.RUnlock()

	for id, register := range match {
		for _, addr := range register {
			log.Printf("send hello: connectionId=%v addr=%v", id, addr)

			_, err := conn.WriteTo([]byte("hello"), addr)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func receiveStun(conn *net.UDPConn) {

	buffer := make([]byte, 1024)

	br, remoteAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatalln(err)
	}

	connectId := string(buffer[0:br])
	log.Printf("receive connectId=%s remoteAddr=%v", connectId, remoteAddr)

	matchMu.Lock()
	defer matchMu.Unlock()

	// Todo: Remove old stuff
	if _, ok := match[connectId]; !ok {
		match[connectId] = make(map[string]*net.UDPAddr)
	}

	timerKey := connectId + "-" + remoteAddr.String()
	if timer, ok := timers[timerKey]; ok {
		timer.Reset(time.Second * 40)
	} else {
		timer := time.NewTimer(time.Second * 40)
		timers[timerKey] = timer

		go func() {
			if _, ok := <-timer.C; ok {

				log.Printf("remove old connection trace connectId=%s", connectId)

				matchMu.Lock()
				defer matchMu.Unlock()

				delete(match[connectId], remoteAddr.String())
				delete(timers, timerKey)

				if len(match[connectId]) == 0 {
					delete(match, connectId)
				}
			}
		}()
	}

	match[connectId][remoteAddr.String()] = remoteAddr

	if len(match[connectId]) != 2 {
		return
	}

	hosts := make([]*net.UDPAddr, 2)

	var inx int
	for _, host := range match[connectId] {
		hosts[inx] = host
		inx++
	}

	_, err = conn.WriteToUDP([]byte(hosts[0].String()), hosts[1])
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.WriteToUDP([]byte(hosts[1].String()), hosts[0])
	if err != nil {
		log.Fatalln(err)
	}

	delete(match, connectId)
}

func Server() {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: stunPort})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("start stun server on %v", conn.LocalAddr())

	go func() {
		for range time.Tick(time.Second * 20) {
			ping(conn)
		}
	}()

	for {
		receiveStun(conn)
	}
}

// GOOS=linux GOARCH=amd64 go build cmd/main.go
// scp main 7zierahn@rzssh1.informatik.uni-hamburg.de:~
