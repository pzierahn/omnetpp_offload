package stargate

import (
	"net"
)

func RelayServerUDP() (addr1, addr2 net.Addr, err error) {
	conn1, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}

	conn2, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}

	addr1 = conn1.LocalAddr()
	addr2 = conn2.LocalAddr()

	log.Printf("RelayServerUDP: addr1=%v addr2=%v", addr1, addr2)

	go func() {
		for {
			// https://stackoverflow.com/questions/2613734/maximum-packet-size-for-a-tcp-connection
			buf := make([]byte, 65535)
			br, err := conn1.Read(buf)
			if err != nil {
				log.Fatalln(err)
			}

			_, err = conn2.Write(buf[:br])
			if err != nil {
				log.Fatalln(err)
			}
		}
	}()

	go func() {
		for {
			// https://stackoverflow.com/questions/2613734/maximum-packet-size-for-a-tcp-connection
			buf := make([]byte, 65535)
			br, err := conn2.Read(buf)
			if err != nil {
				log.Fatalln(err)
			}

			_, err = conn1.Write(buf[:br])
			if err != nil {
				log.Fatalln(err)
			}
		}
	}()

	return
}

func RelayDialUDP(addr net.Addr) (conn *net.UDPConn, err error) {

	log.Printf("RelayDialUDP: addr=%v", addr)

	laddr := &net.UDPAddr{}
	raddr, err := net.ResolveUDPAddr("udp", addr.String())
	conn, err = net.DialUDP("tcp", laddr, raddr)

	return
}
