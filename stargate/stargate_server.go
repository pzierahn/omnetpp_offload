package stargate

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type clientType map[string]bool

var clients = clientType{}

func (c clientType) keys(filter string) string {
	output := []string{}
	for key := range c {
		if key != filter {
			output = append(output, key)
		}
	}

	return strings.Join(output, ",")
}

// Server --
func Server() {
	localAddress := ":9595"
	if len(os.Args) > 2 {
		localAddress = os.Args[2]
	}

	addr, _ := net.ResolveUDPAddr("udp", localAddress)
	conn, _ := net.ListenUDP("udp", addr)

	fmt.Printf("addr=%v localAddress=%v\n", addr, localAddress)

	for {
		buffer := make([]byte, 1024)
		bytesRead, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}

		incoming := string(buffer[0:bytesRead])
		fmt.Println("[INCOMING]", incoming, remoteAddr)
		if incoming != "register" {
			continue
		}

		//_, err = conn.WriteToUDP([]byte("Ping from sever"), remoteAddr)
		//if err != nil {
		//	panic(err)
		//}

		clients[remoteAddr.String()] = true

		for client := range clients {
			resp := clients.keys(client)
			if len(resp) > 0 {
				r, err := net.ResolveUDPAddr("udp", client)
				if err != nil {
					panic(err)
				}

				_, err = conn.WriteTo([]byte(resp), r)
				if err != nil {
					panic(err)
				}

				fmt.Printf("[INFO] Responded to %s with %s\n", client, string(resp))
			}
		}
	}
}

// GOOS=linux GOARCH=amd64 go build cmd/main.go
// scp main 7zierahn@rzssh1.informatik.uni-hamburg.de:~
