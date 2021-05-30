package stargate

import (
	"log"
	"net"
)

func Server() {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: stunPort})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("stun server on %v", conn.LocalAddr())

	cache := make(map[string][]string)
	buffer := make([]byte, 1024)

	for {
		bytesRead, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatalln(err)
		}

		connectId := string(buffer[0:bytesRead])
		log.Printf("connectId %s %v", connectId, remoteAddr)

		cache[connectId] = append(cache[connectId], remoteAddr.String())

		if len(cache[connectId]) != 2 {
			continue
		}

		host1, _ := net.ResolveUDPAddr("udp", cache[connectId][0])
		host2, _ := net.ResolveUDPAddr("udp", cache[connectId][1])

		_, err = conn.WriteToUDP([]byte(host2.String()), host1)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = conn.WriteToUDP([]byte(host1.String()), host2)
		if err != nil {
			log.Fatalln(err)
		}

		delete(cache, connectId)
	}
}

// GOOS=linux GOARCH=amd64 go build cmd/main.go
// scp main 7zierahn@rzssh1.informatik.uni-hamburg.de:~
