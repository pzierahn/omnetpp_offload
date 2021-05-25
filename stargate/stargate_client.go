package stargate

import (
	"github.com/patrickz98/project.go.omnetpp/quick"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"log"
	"net"
	"os"
	"time"
)

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

	qConn := &quick.Connection{
		Connection: conn,
	}

	qConn.Init()

	log.Printf("remote=%v local=%v", remote, local)

	//go func() {
	//	time.Sleep(time.Second)

	log.Printf("sending register")

	registerMsg := "register"
	err := qConn.Send(registerMsg, remote)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("register done")
	//}()

	//listen(qConn)
}

func listen(conn *quick.Connection) {
	//buffer := make([]byte, 1024*1024)

	for {
		var garbage string
		remoteAddr, err := conn.Receive(&garbage)
		if err != nil {
			log.Println("[ERROR]", err)
			continue
		}

		log.Printf("received remote=%v garbage=%d", remoteAddr, len(garbage))

		//var addrs []*net.UDPAddr
		//err := qConn.Receive(&addrs)
		//if err != nil {
		//	log.Println("[ERROR]", err)
		//	continue
		//}
		//
		//log.Printf("received addrs=%s", simple.PrettyString(addrs))

		////for _, a := range strings.Split(string(buffer[0:bytesRead]), ",") {
		////	if a != local {
		////		go chatter(conn, a)
		////	}
		////}
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
