package stargate

import (
	"context"
	"encoding/json"
	"net"
)

func PropagateTCP(ctx context.Context, dialAddr DialAddr, addr *net.TCPAddr) (err error) {

	log.Printf("PropagateTCP: dialAddr=%v addr=%v", dialAddr, addr)

	conn, err := net.ListenMulticastUDP("udp", nil, broadcast)
	if err != nil {
		return
	}

	defer func() { _ = conn.Close() }()

	byt, err := json.MarshalIndent(addr, "", "  ")
	if err != nil {
		return
	}

	go func() {
		for {
			buf := make([]byte, 1024)
			br, raddr, err := conn.ReadFrom(buf)
			if err != nil {
				break
			}

			requestAddr := string(buf[:br])
			log.Printf("PropagateTCP: requested dialAddr %v from %v", requestAddr, raddr)

			if dialAddr != requestAddr {
				continue
			}

			log.Printf("PropagateTCP: write %v", string(byt))

			_, err = conn.WriteTo(byt, raddr)
			if err != nil {
				break
			}
		}
	}()

	select {
	case <-ctx.Done():
	}

	return
}

func DialLocal(ctx context.Context, dialAddr DialAddr) (raddr net.TCPAddr, err error) {
	bc, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}

	if deadline, ok := ctx.Deadline(); ok {
		err = bc.SetDeadline(deadline)
		if err != nil {
			return
		}
	}

	_, err = bc.WriteToUDP([]byte(dialAddr), broadcast)
	if err != nil {
		return
	}

	buf := make([]byte, 1024)

	br, uaddr, err := bc.ReadFromUDP(buf)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf[:br], &raddr)
	if err != nil {
		return
	}

	raddr.IP = uaddr.IP
	raddr.Zone = uaddr.Zone

	return
}
