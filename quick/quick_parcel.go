package quick

import (
	"fmt"
	"net"
)

type parcelList []bool

type cast int

const (
	castParcel cast = iota + 1
	castPairing
	castPing
	castList
)

type parcel struct {
	sessionId uint32
	messageId uint32
	cast      cast
	offset    uint64
	payload   []byte
}

func (conn *Connection) sendParcel(parcel *parcel, addr *net.UDPAddr) (err error) {

	var buf []byte
	buf, err = encodeGob(parcel)
	if err != nil {
		err = fmt.Errorf("error: gobbing parcel: %v", err)
		return
	}

	conn.connMu.Lock()
	_, err = conn.Connection.WriteToUDP(buf, addr)
	conn.connMu.Unlock()

	if err != nil {
		err = fmt.Errorf("error: sending parcel: %v", err)
		return
	}

	return
}
