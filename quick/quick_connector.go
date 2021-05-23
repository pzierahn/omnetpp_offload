package quick

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"log"
	"math"
	"math/rand"
	"net"
	"sync"
)

const (
	payloadSize = 1024
)

type ParcelType int

const (
	TypeParcel ParcelType = iota + 1
	TypeAck
	TypeConfirm
)

type Connection struct {
	sync.Mutex
	Connection *net.UDPConn
}

func parcelId(prefix, suffix uint32) (id int64) {
	id = int64(prefix) << 32
	id |= int64(suffix)

	return
}

func (conn *Connection) Send(obj interface{}, addr *net.UDPAddr) (err error) {
	conn.Lock()
	defer conn.Unlock()

	//
	// Encode struct to gob bytes
	//

	var objBuf bytes.Buffer
	objEnc := gob.NewEncoder(&objBuf)
	err = objEnc.Encode(obj)

	if err != nil {
		err = fmt.Errorf("error gobbing interface %T: %v", obj, err)
		return
	}

	//
	// Compile and send parcels
	//

	id := rand.Uint32()
	chunks := uint32(math.Ceil(float64(objBuf.Len()) / float64(payloadSize)))

	log.Printf("parcelId=%x size=%d chunks=%d", id, objBuf.Len(), chunks)

	for idx := uint32(0); idx < chunks; idx++ {

		endSlice := simple.MathMin(int(payloadSize*(idx+1)), objBuf.Len())

		parcel := &Parcel{
			Type:    TypeParcel,
			Package: id,
			Index:   idx,
			Chunks:  chunks,
			Payload: objBuf.Bytes()[payloadSize*idx : endSlice],
		}

		log.Printf("sending parcel=%d payload='%s' (%d bytes)",
			idx, objBuf.Bytes()[payloadSize*idx:endSlice], len(parcel.Payload))

		var buf []byte
		buf, err = parcel.marshalGob()
		if err != nil {
			err = fmt.Errorf("error gobbing parcel: %v", err)
			return
		}

		_, err = conn.Connection.WriteToUDP(buf, addr)
		if err != nil {
			err = fmt.Errorf("error sending parcel: %v", err)
			return
		}
	}

	return
}

func (conn *Connection) Receive(obj interface{}) (err error) {

	buffer := make([]byte, 1024+payloadSize)

	var message []byte
	var size int

	var parcels = -1
	log.Printf("listening on %v", conn.Connection.LocalAddr())

	for parcels != 0 {
		var bytesRead int
		bytesRead, err = conn.Connection.Read(buffer)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		var parcel Parcel

		dec := gob.NewDecoder(bytes.NewReader(buffer[0:bytesRead]))
		err = dec.Decode(&parcel)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		log.Printf("recieved: index=%d payload='%s'", parcel.Index, parcel.Payload)

		if message == nil {
			log.Printf("recieved: first parcel Chunks=%v", parcel.Chunks)

			parcels = int(parcel.Chunks)
			message = make([]byte, payloadSize*parcel.Chunks)
		}

		copy(message[payloadSize*parcel.Index:payloadSize*(parcel.Index+1)], parcel.Payload)
		size += len(parcel.Payload)

		parcels--
	}

	log.Printf("size=%d", size)
	//log.Printf("message: '%s'", message[0:size])

	enc := gob.NewDecoder(bytes.NewReader(message))
	err = enc.Decode(obj)

	return
}
