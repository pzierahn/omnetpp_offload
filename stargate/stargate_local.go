package stargate

import (
	"context"
	"encoding/json"
	"net"
	"time"
)

var broadcast = &net.UDPAddr{
	IP:   net.IPv4(239, 11, 22, 33),
	Port: 10077,
}

// BroadcastTCP will listen for multicast broadcasts. It will respond with the addr if the dialAddr matches.
func BroadcastTCP(ctx context.Context, dialAddr DialAddr, addr *net.TCPAddr) (err error) {

	log.Printf("BroadcastTCP: dialAddr=%v addr=%v", dialAddr, addr)

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
			if dialAddr != requestAddr {
				continue
			}

			log.Printf("BroadcastTCP: requestAddr=%v responseTo=%v write=%v",
				requestAddr, raddr, string(byt))

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

// DialLocal will broadcast the dialAddr to the local network.
// It returns a TCP address on which to connect to peers.
func DialLocal(ctx context.Context, dialAddr DialAddr) (raddr net.TCPAddr, err error) {

	log.Printf("DialLocal: dialAddr=%v", dialAddr)

	bc, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}

	ctx, cnl := context.WithTimeout(ctx, time.Millisecond*1000)
	defer cnl()

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

	// TODO: Shrink buffer size
	buf := make([]byte, 1024)

	br, origin, err := bc.ReadFromUDP(buf)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf[:br], &raddr)
	if err != nil {
		return
	}

	raddr.IP = origin.IP
	raddr.Zone = origin.Zone

	return
}
