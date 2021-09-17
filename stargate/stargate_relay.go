package stargate

import (
	"context"
	"fmt"
	"net"
	"time"
)

const (
	success = "ok"
)

func (server *stargateServer) relayTCPServer() (err error) {

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: config.Port,
	})
	if err != nil {
		return
	}

	log.Printf("start stargate relay server on %v", lis.Addr())

	for {
		conn, err := lis.AcceptTCP()
		if err != nil {
			log.Fatalln(err)
		}

		go server.rendezvousTCP(conn)
	}
}

func (server *stargateServer) rendezvousTCP(conn *net.TCPConn) {
	buf := make([]byte, 1024)
	br, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}

	dialAddr := string(buf[:br])

	log.Printf("rendezvousTCP: dialAddr='%s' LocalAddr=%v RemoteAddr=%v",
		dialAddr, conn.LocalAddr(), conn.RemoteAddr())

	server.relayMu.Lock()
	defer server.relayMu.Unlock()

	peer, ok := server.relay[dialAddr]

	if !ok {

		//
		// Wait for peer...
		//

		server.relay[dialAddr] = conn
		return
	} else {

		//
		// Connect both peers
		//

		delete(server.relay, dialAddr)

		_, err = peer.Write([]byte(success))
		if err != nil {
			log.Println(err)
			return
		}

		_, err = conn.Write([]byte(success))
		if err != nil {
			log.Println(err)
			return
		}

		pipeAllTCP(peer, conn)
	}
}

func pipeTCP(from, to *net.TCPConn) {
	defer func() { _, _ = from.Close(), to.Close() }()

	for {
		// https://stackoverflow.com/questions/2613734/maximum-packet-size-for-a-tcp-connection
		buf := make([]byte, 65535)
		br, err := from.Read(buf)
		if err != nil {
			//log.Println(err)
			break
		}

		_, err = to.Write(buf[:br])
		if err != nil {
			//log.Println(err)
			break
		}
	}
}

func pipeAllTCP(conn1, conn2 *net.TCPConn) {
	go pipeTCP(conn1, conn2)
	go pipeTCP(conn2, conn1)
}

// DialRelayTCP will establish a TCP relay connection over the stargate server.
func DialRelayTCP(ctx context.Context, dial DialAddr) (conn *net.TCPConn, err error) {

	raddr, err := config.tcpAddr()
	if err != nil {
		return
	}

	log.Printf("DialRelayTCP: dial=%v addr=%v", dial, raddr)

	conn, err = net.DialTCP("tcp", &net.TCPAddr{}, raddr)
	if err != nil {
		return
	}

	// Set connection timeout
	if deadline, ok := ctx.Deadline(); ok {
		err = conn.SetDeadline(deadline)
		if err != nil {
			return
		}

		// Reset deadline
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

	if string(buf[:br]) != success {
		err = fmt.Errorf("connection failed: read '%v'", string(buf[:br]))
	}

	return
}
