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
	//payloadSize = 8192
	payloadSize = 1
)

type Connection struct {
	sync.Mutex
	Connection    *net.UDPConn
	parcelListRec map[uint32]chan *ParcelWithAddress
	ping          chan *ParcelWithAddress
	receiver      chan *ParcelWithAddress
	connMu        sync.Mutex
}

func (conn *Connection) Init() {
	conn.parcelListRec = make(map[uint32]chan *ParcelWithAddress)

	go func() {
		buffer := make([]byte, 65536)

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

			if pkg.Type == TypeParcelList {
				log.Printf("receive list: %x", pkg.MessageId)

				if ackList, ok := conn.parcelListRec[pkg.MessageId]; ok {
					ackList <- &pkg
				}
			}

			if pkg.Type == TypeParcel {
				log.Printf("receive parcel: %x", pkg.MessageId)
				if conn.receiver != nil {
					conn.receiver <- &pkg
				}
			}

			if pkg.Type == TypePing {
				log.Printf("receive ping: %x", pkg.MessageId)
				if conn.ping != nil {
					conn.ping <- &pkg
				}
			}

			conn.Unlock()
		}
	}()
}

func (conn *Connection) sendParcel(parcel *Parcel, addr *net.UDPAddr) (err error) {

	var buf []byte
	buf, err = parcel.marshalGob()
	if err != nil {
		err = fmt.Errorf("error gobbing parcel: %v", err)
		return
	}

	conn.connMu.Lock()
	defer conn.connMu.Unlock()

	_, err = conn.Connection.WriteToUDP(buf, addr)
	if err != nil {
		err = fmt.Errorf("error sending parcel: %v", err)
		return
	}

	return
}

func (conn *Connection) Send(obj interface{}, addr *net.UDPAddr) (err error) {

	//
	// Prepare message
	//

	var payload []byte
	payload, err = encode(obj)
	if err != nil {
		return
	}

	messageId := rand.Uint32()
	chunks := uint32(math.Ceil(float64(len(payload)) / float64(payloadSize)))
	parcels := make([]*Parcel, chunks)

	for idx := uint32(0); idx < chunks; idx++ {

		endSlice := simple.MathMin(int(payloadSize*(idx+1)), len(payload))

		parcel := &Parcel{
			Type:      TypeParcel,
			MessageId: messageId,
			Index:     idx,
			Chunks:    chunks,
			Payload:   payload[payloadSize*idx : endSlice],
		}

		parcels[idx] = parcel
	}

	log.Printf("create new message: id=%x size=%d chunks=%d", messageId, len(payload), chunks)

	//
	// TypeParcelList receiver
	//

	done := make(chan bool)
	sender := make(chan *Parcel, 8)

	parcelLists := make(chan *ParcelWithAddress)
	conn.Lock()
	conn.parcelListRec[messageId] = parcelLists
	conn.Unlock()

	defer func() {
		conn.Lock()
		delete(conn.parcelListRec, messageId)
		conn.Unlock()
		close(parcelLists)
	}()

	go func() {
		for list := range parcelLists {

			var receivedPackages parcelList
			err = decode(list.Payload, &receivedPackages)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			log.Printf("recieved list %v: %v",
				list.MessageId, simple.PrettyString(receivedPackages))

			var missingParcels int

			for parcelId, val := range receivedPackages {
				if val {
					continue
				}

				sender <- parcels[parcelId]
				missingParcels++
			}

			if missingParcels == 0 {
				//
				// Message delivered! Exit send function now!
				//

				done <- true
				break
			}
		}

		log.Printf("parcel lister finished id=%x", messageId)
	}()

	//
	// TypePing sender
	//

	pingTic := time.NewTimer(time.Millisecond * 50)

	go func() {
		for range pingTic.C {
			log.Printf("Send: ping message id=%x chunks=%d", messageId, chunks)

			ack := &Parcel{
				Type:      TypePing,
				MessageId: messageId,
				Chunks:    chunks,
			}

			err = conn.sendParcel(ack, addr)
			if err != nil {
				err = fmt.Errorf("error sending parcel: %v", err)
				break
			}

			pingTic.Reset(time.Millisecond * 50)
		}

		log.Printf("ack message finished id=%x", messageId)
	}()

	//
	// Send parcels
	//

	go func() {
		for parcel := range sender {
			log.Printf("Send: sender: sending id=%x idx=%d (%d bytes)",
				messageId, parcel.Index, len(parcel.Payload))

			err = conn.sendParcel(parcel, addr)
			if err != nil {
				log.Printf("error: %v", err)
				break
			}
		}

		log.Printf("Send: sender: finished id=%x", messageId)
	}()

	<-done

	log.Printf("Send: finishing id=%x", messageId)
	pingTic.Stop()

	//select {
	//case <-done:
	//case <-time.After(2 * time.Second):
	//	log.Printf("Send: timeout id=%x", messageId)
	//}

	return
}

func (conn *Connection) Receive(obj interface{}) (remoteAddr *net.UDPAddr, err error) {

	log.Printf("listening on %v", conn.Connection.LocalAddr())

	var isInit bool
	var messageId uint32
	var chunks uint32

	var mu sync.RWMutex
	var size int
	var message []byte
	var received parcelList
	var receivedParcels uint32

	var wg sync.WaitGroup

	init := func(ping *ParcelWithAddress) {
		messageId = ping.MessageId
		chunks = ping.Chunks
		received = make(parcelList, ping.Chunks)
		message = make([]byte, ping.Chunks*payloadSize)
		remoteAddr = ping.RemoteAddr
		isInit = true
	}

	//
	// Ping receiver
	//

	killParcelListSender := make(chan bool)
	defer close(killParcelListSender)

	pings := make(chan *ParcelWithAddress)

	conn.Lock()
	conn.ping = pings
	conn.Unlock()

	wg.Add(1)
	go func() {
		timeout := time.NewTimer(time.Second * 10)

	loop:
		for {
			select {
			case ping := <-pings:
				if !isInit {
					init(ping)
				}

				log.Printf("Receive: ping messageId=%x chunks=%v", ping.MessageId, ping.Chunks)
				timeout.Reset(time.Millisecond * 500)

			case <-timeout.C:
				log.Printf("Receive: timeout messageId=%x", messageId)
				break loop
			}
		}

		log.Printf("Receive: ping finished")
		//killParcelListSender <- true
		wg.Done()
	}()

	//
	// Received parcels sender
	//

	listTic := time.NewTimer(time.Millisecond * 75)

	wg.Add(1)
	go func() {

	loop:
		for {

			select {
			case <-listTic.C:
				if !isInit {
					listTic.Reset(time.Millisecond * 75)
					continue
				}

				var payload []byte
				mu.RLock()
				log.Printf("Receive: send parcel list messageId=%x", messageId)
				//messageId, simple.PrettyString(received))
				payload, err = encode(received)
				mu.RUnlock()

				list := &Parcel{
					Type:      TypeParcelList,
					MessageId: messageId,
					Payload:   payload,
				}

				err = conn.sendParcel(list, remoteAddr)
				if err != nil {
					log.Fatalln(err)
				}

				listTic.Reset(time.Millisecond * 75)
			case <-killParcelListSender:
				break loop
			}
		}

		log.Printf("Receive: parcellist sender finished")
		wg.Done()
	}()

	//
	// Parcel receiver
	//

	parcels := make(chan *ParcelWithAddress, 8)

	conn.Lock()
	conn.receiver = parcels
	conn.Unlock()

	for parcel := range parcels {

		if !isInit {
			continue
		}

		log.Printf("Receive: recieved: index=%4d (%d bytes)", parcel.Index, len(parcel.Payload))

		windowStart := payloadSize * parcel.Index
		windowEnd := windowStart + uint32(len(parcel.Payload))
		copy(message[windowStart:windowEnd], parcel.Payload)
		size += len(parcel.Payload)

		var receivedAll bool
		mu.Lock()
		received[parcel.Index] = true
		receivedParcels++
		receivedAll = receivedParcels == chunks
		mu.Unlock()

		if receivedAll {
			break
		}
	}

	log.Printf("Receive: recieved all parcels size=%d", size)

	//listTic.Reset(0)
	//listTic.Stop()

	wg.Wait()

	//log.Printf("message: '%s'", message[0:size])

	enc := gob.NewDecoder(bytes.NewReader(message))
	err = enc.Decode(obj)

	return
}
