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
	"time"
)

const (
	payloadSize = 8192
	ackWait     = 100 * time.Millisecond
)

type ParcelType int

const (
	TypeParcel ParcelType = iota + 1
	TypeAck
)

type Connection struct {
	sync.Mutex
	Connection  *net.UDPConn
	ackListener map[uint32]chan *ParcelWithAddress
	receiver    chan *ParcelWithAddress
}

func (conn *Connection) Init() {
	conn.ackListener = make(map[uint32]chan *ParcelWithAddress)

	go func() {
		buffer := make([]byte, 1024+payloadSize)

		for {
			bytesRead, remote, err := conn.Connection.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			var parcel Parcel
			err = parcel.unmarshalGob(buffer[0:bytesRead])
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			pkg := ParcelWithAddress{
				Parcel:     &parcel,
				RemoteAddr: remote,
			}

			conn.Lock()

			log.Printf("recieved message: %x %v", pkg.Id, pkg.Type)

			if pkg.Type == TypeAck {
				if ackList, ok := conn.ackListener[pkg.Id]; ok {
					ackList <- &pkg
				}
			}

			if pkg.Type == TypeParcel {
				if conn.receiver != nil {
					conn.receiver <- &pkg
				}
			}

			conn.Unlock()
		}
	}()
}

func (conn *Connection) Send(obj interface{}, addr *net.UDPAddr) (err error) {

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
	// Setup
	//

	done := make(chan bool)
	defer close(done)

	id := rand.Uint32()
	chunks := uint32(math.Ceil(float64(objBuf.Len()) / float64(payloadSize)))
	cache := make([]*Parcel, chunks)
	acks := make(chan *ParcelWithAddress)

	conn.Lock()
	conn.ackListener[id] = acks
	conn.Unlock()

	log.Printf("id=%x size=%d chunks=%d", id, objBuf.Len(), chunks)

	//
	// Receive ack messages
	//

	go func() {
		for parcel := range acks {

			var receivedPackages map[uint32]bool
			dec := gob.NewDecoder(bytes.NewReader(parcel.Payload))
			err = dec.Decode(&receivedPackages)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			log.Printf("recieved ack %v: %v", parcel.Id, simple.PrettyString(receivedPackages))

			if len(cache) == len(receivedPackages) {
				log.Printf("message deliverd successful!")

				conn.Lock()
				delete(conn.ackListener, id)
				close(acks)
				conn.Unlock()

				done <- true
			}
		}
	}()

	//
	// Compile and send parcels
	//

	for idx := uint32(0); idx < chunks; idx++ {

		endSlice := simple.MathMin(int(payloadSize*(idx+1)), objBuf.Len())

		parcel := &Parcel{
			Type:    TypeParcel,
			Id:      id,
			Index:   idx,
			Chunks:  chunks,
			Payload: objBuf.Bytes()[payloadSize*idx : endSlice],
		}

		cache[idx] = parcel

		log.Printf("sending parcel=%d (%d bytes)",
			idx, len(parcel.Payload))

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

	<-done

	log.Printf("parcelId=%x send all packages", id)

	return
}

func (conn *Connection) Receive(obj interface{}) (remoteAddr *net.UDPAddr, err error) {

	var messageId uint32

	var message []byte
	var size int

	var parcels = -1
	log.Printf("listening on %v", conn.Connection.LocalAddr())

	var mu sync.Mutex
	received := make(map[uint32]bool)

	done := make(chan bool)
	defer close(done)

	receiver := make(chan *ParcelWithAddress)

	conn.Lock()
	conn.receiver = receiver
	conn.Unlock()

	go func() {
		time.Sleep(time.Millisecond * 100)

	loop:
		for {
			select {
			case <-time.Tick(ackWait):
				if messageId == 0 {
					continue
				}

				if messageId == 0 {
					continue
				}

				mu.Lock()

				log.Printf("send ack %x: %v", messageId, simple.PrettyString(received))

				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				err = enc.Encode(received)
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}

				mu.Unlock()

				ack := Parcel{
					Type:    TypeAck,
					Id:      messageId,
					Payload: buf.Bytes(),
				}

				byt, err := ack.marshalGob()
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}

				_, err = conn.Connection.WriteToUDP(byt, remoteAddr)
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}

			case <-done:
				log.Printf("ack sender: quits")

				if messageId == 0 {
					continue
				}

				if messageId == 0 {
					continue
				}

				mu.Lock()

				log.Printf("send ack %x: %v", messageId, simple.PrettyString(received))

				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				err = enc.Encode(received)
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}

				mu.Unlock()

				ack := Parcel{
					Type:    TypeAck,
					Id:      messageId,
					Payload: buf.Bytes(),
				}

				byt, err := ack.marshalGob()
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}

				_, err = conn.Connection.WriteToUDP(byt, remoteAddr)
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}

				break loop
			}
		}
	}()

	for parcel := range receiver {

		go func() {
			mu.Lock()
			received[parcel.Index] = true
			mu.Unlock()
		}()

		log.Printf("recieved: index=%4d (%d bytes)", parcel.Index, len(parcel.Payload))

		if messageId == 0 {
			log.Printf("recieved: init parcel")

			parcels = int(parcel.Chunks)
			messageId = parcel.Id
			remoteAddr = parcel.RemoteAddr
			message = make([]byte, payloadSize*parcel.Chunks)

			ackStream := make(chan *ParcelWithAddress)
			conn.Lock()
			conn.ackListener[messageId] = ackStream
			conn.Unlock()

			go func() {
				for ack := range ackStream {
					var list map[int]bool

					dec := gob.NewDecoder(bytes.NewReader(ack.Payload))
					err = dec.Decode(&list)
					if err != nil {
						log.Printf("error: %v", err)
						continue
					}

					log.Printf("recived ack: %v", simple.PrettyString(list))

					if len(list) == parcels {
						done <- true
						break
					}
				}
			}()
		}

		sliceStart := payloadSize * parcel.Index
		copy(message[sliceStart:sliceStart+uint32(len(parcel.Payload))], parcel.Payload)
		size += len(parcel.Payload)

		parcels--

		if parcels == 0 {

			log.Printf("recieved all parcels: closing stuff")

			//done <- true

			go func() {
				conn.Lock()
				conn.receiver = nil
				close(receiver)
				conn.Unlock()

				log.Printf("recieved all parcels: closing done")
			}()

			//
			//break
		}
	}

	log.Printf("size=%d", size)
	//log.Printf("message: '%s'", message[0:size])

	enc := gob.NewDecoder(bytes.NewReader(message))
	err = enc.Decode(obj)

	return
}
