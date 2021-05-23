package stargate

import (
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/quick"
	"log"
	"net"
	"os"
	"strings"
)

type clientType map[string]bool

var clients = clientType{}

func (c clientType) keys(filter string) string {
	var output []string
	for key := range c {
		if key != filter {
			output = append(output, key)
		}
	}

	return strings.Join(output, ",")
}

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

		clients[strings.TrimSpace(remoteAddr.String())] = true

		for client := range clients {
			resp := clients.keys(client)
			if len(resp) > 0 {
				var remote *net.UDPAddr
				remote, err = net.ResolveUDPAddr("udp", client)
				if err != nil {
					log.Fatalln(err)
				}

				quickConn := quick.Connection{
					Connection: conn,
				}

				log.Printf("resp: %v", resp)

				err = quickConn.Send(resp, remote)
				if err != nil {
					log.Fatalln(err)
				}

				//var buffer bytes.Buffer
				//enc := gob.NewEncoder(&buffer)
				//err = enc.Encode(quick.Parcel{
				//	Payload: []byte(resp),
				//})
				//if err != nil {
				//	log.Fatalln(err)
				//}
				//
				//_, err = conn.WriteTo(buffer.Bytes(), remote)
				//if err != nil {
				//	log.Fatalln(err)
				//}

				log.Printf("[INFO] Responded to %s with %s", client, resp)
			}
		}
	}
}

// GOOS=linux GOARCH=amd64 go build cmd/main.go
// scp main 7zierahn@rzssh1.informatik.uni-hamburg.de:~
