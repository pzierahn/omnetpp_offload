package stargate

import (
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Parcel struct {
	Package uint32
	Index   uint32
	Chunks  uint32
	Payload []byte
}

func Client() {
	register()
}

func register() {
	signalAddress := os.Args[2]

	localAddress := ":9595"
	if len(os.Args) > 3 {
		localAddress = os.Args[3]
	}

	remote, _ := net.ResolveUDPAddr("udp", signalAddress)
	local, _ := net.ResolveUDPAddr("udp", localAddress)
	conn, _ := net.ListenUDP("udp", local)

	log.Printf("remote=%v local=%v", remote, local)

	go func() {
		time.Sleep(time.Second)

		bytesWritten, err := conn.WriteTo([]byte("register"), remote)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("register (%v bytes)", bytesWritten)
	}()

	listen(conn, local.String())
}

func listen(conn *net.UDPConn, local string) {
	for {
		log.Printf("listening on %v", conn.LocalAddr())
		buffer := make([]byte, 1024)
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}

		log.Printf("recieved: '%s'", buffer[0:bytesRead])
		if string(buffer[0:bytesRead]) == "Hello!" {
			continue
		}

		for _, a := range strings.Split(string(buffer[0:bytesRead]), ",") {
			if a != local {
				go chatter(conn, a)
			}
		}
	}
}

func chatter(conn *net.UDPConn, remote string) {
	addr, _ := net.ResolveUDPAddr("udp", remote)
	for {
		message := simple.NamedId("message", 4)
		_, err := conn.WriteTo([]byte(message), addr)
		if err != nil {
			//log.Fatalln("[ERROR]", err)
			continue
		}

		log.Printf("sent: '%s' to %v", message, remote)
		time.Sleep(5 * time.Second)
	}
}
