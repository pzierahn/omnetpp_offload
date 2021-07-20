package stargate

import (
	"net"
)

func RelayServerTCP() (port1, port2 int, err error) {

	listener1, err := net.Listen("tcp", ":0")
	if err != nil {
		return
	}

	listener2, err := net.Listen("tcp", ":0")
	if err != nil {
		return
	}

	port1 = listener1.Addr().(*net.TCPAddr).Port
	port2 = listener2.Addr().(*net.TCPAddr).Port

	log.Printf("RelayServerTCP: port1=%v port2=%v", port1, port2)

	incoming := make(chan net.Conn)

	go func() {
		conn, err := listener1.Accept()
		if err != nil {
			log.Println(err)
			incoming <- nil
			return
		}

		log.Printf("RelayServerTCP: LocalAddr=%v RemoteAddr=%v", conn.LocalAddr(), conn.RemoteAddr())
		incoming <- conn
	}()

	go func() {
		conn, err := listener2.Accept()
		if err != nil {
			log.Println(err)
			incoming <- nil
			return
		}

		log.Printf("RelayServerTCP: LocalAddr=%v RemoteAddr=%v", conn.LocalAddr(), conn.RemoteAddr())
		incoming <- conn
	}()

	go func() {
		conn1 := <-incoming
		conn2 := <-incoming
		close(incoming)

		if conn1 == nil || conn2 == nil {
			return
		}

		go func() {
			for {
				// https://stackoverflow.com/questions/2613734/maximum-packet-size-for-a-tcp-connection
				buf := make([]byte, 65535)
				br, err := conn1.Read(buf)
				if err != nil {
					log.Println(err)
					break
				}

				_, err = conn2.Write(buf[:br])
				if err != nil {
					log.Println(err)
					break
				}
			}

			_ = conn1.Close()
		}()

		go func() {
			for {
				// https://stackoverflow.com/questions/2613734/maximum-packet-size-for-a-tcp-connection
				buf := make([]byte, 65535)
				br, err := conn2.Read(buf)
				if err != nil {
					log.Println(err)
					break
				}

				_, err = conn1.Write(buf[:br])
				if err != nil {
					log.Println(err)
					break
				}
			}

			_ = conn2.Close()
		}()
	}()

	return
}

func RelayDialTCP(addr net.Addr) (conn *net.TCPConn, err error) {

	log.Printf("RelayDialTCP: addr=%v", addr)

	laddr := &net.TCPAddr{}
	tcpaddr, err := net.ResolveTCPAddr("tcp", addr.String())
	conn, err = net.DialTCP("tcp", laddr, tcpaddr)

	return
}
