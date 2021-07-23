package stargate

import (
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"net"
	"sync"
	"time"
)

const (
	relaySuccessful = "ok"
)

var dialMu sync.Mutex
var relay = make(map[DialAddr]*net.TCPConn)

//
// TODO: Close connections!
//

func ServerRelayTCP() (err error) {

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: config.Port,
	})
	if err != nil {
		return
	}

	log.Printf("ServerRelayTCP: started on %v", lis.Addr())

	for {
		conn, err := lis.AcceptTCP()
		if err != nil {
			log.Fatalln(err)
		}

		go rendezvousTCP(conn)
	}
}

func rendezvousTCP(conn *net.TCPConn) {
	buf := make([]byte, 1024)
	br, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}

	dialAddr := string(buf[:br])

	log.Printf("rendezvousTCP: dialAddr='%s' LocalAddr=%v RemoteAddr=%v",
		dialAddr, conn.LocalAddr(), conn.RemoteAddr())

	dialMu.Lock()
	defer dialMu.Unlock()

	peer, ok := relay[dialAddr]

	if !ok {
		relay[dialAddr] = conn
		return
	}

	delete(relay, dialAddr)

	_, err = peer.Write([]byte(relaySuccessful))
	if err != nil {
		log.Println(err)
		return
	}

	_, err = conn.Write([]byte(relaySuccessful))
	if err != nil {
		log.Println(err)
		return
	}

	pipeAllTCP(peer, conn)
}

func pipeTCP(from, to *net.TCPConn) {
	for {
		// https://stackoverflow.com/questions/2613734/maximum-packet-size-for-a-tcp-connection
		buf := make([]byte, 65535)
		br, err := from.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}

		_, err = to.Write(buf[:br])
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func pipeAllTCP(conn1, conn2 *net.TCPConn) {
	go pipeTCP(conn1, conn2)
	go pipeTCP(conn2, conn1)
}

func RelayDialTCP(ctx context.Context, dial DialAddr) (conn *net.TCPConn, err error) {

	addr := gconfig.StargateDialAddr()
	log.Printf("RelayDialTCP: dial=%v addr=%v", dial, addr)

	laddr := &net.TCPAddr{}
	tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}

	conn, err = net.DialTCP("tcp", laddr, tcpaddr)
	if err != nil {
		return
	}

	if deadline, ok := ctx.Deadline(); ok {
		err = conn.SetDeadline(deadline)
		if err != nil {
			return
		}

		defer func() {
			_ = conn.SetDeadline(time.Time{})
		}()
	}

	_, err = conn.Write([]byte(dial))
	if err != nil {
		return
	}

	buf := make([]byte, 1024)
	br, err := conn.Read(buf)
	if err != nil {
		return
	}

	if string(buf[:br]) != relaySuccessful {
		err = fmt.Errorf("connection failed: wrong relaySuccessful message '%v'", string(buf[:br]))
	}

	return
}
