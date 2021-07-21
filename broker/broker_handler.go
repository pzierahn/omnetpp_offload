package broker

import (
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"log"
	"math/rand"
	"sync"
)

type providerId = string

type broker struct {
	pb.UnimplementedBrokerServer
	pmu         sync.RWMutex
	providers   map[providerId]*pb.ProviderInfo
	umu         sync.RWMutex
	utilization map[providerId]*pb.Utilization
	lmu         sync.RWMutex
	listener    map[string]chan<- *pb.Providers
}

func (broker *broker) providerList() (providers *pb.Providers) {

	providers = &pb.Providers{}

	broker.pmu.RLock()
	for _, prov := range broker.providers {
		providers.Items = append(providers.Items, prov)
	}
	broker.pmu.RUnlock()

	return
}

// GetProviders sends a provider list to the consumer. With every list update an event will be dispatched.
func (broker *broker) GetProviders(_ *pb.Empty, stream pb.Broker_GetProvidersServer) (err error) {

	log.Printf("GetProviders:")

	ctx := stream.Context()

	id := fmt.Sprintf("%x", rand.Uint32())
	listener := make(chan *pb.Providers)
	defer close(listener)

	broker.lmu.Lock()
	broker.listener[id] = listener
	broker.lmu.Unlock()

	defer func() {
		broker.lmu.Lock()
		delete(broker.listener, id)
		broker.lmu.Unlock()
	}()

	//
	// Send initial provider list
	//

	err = stream.Send(broker.providerList())
	if err != nil {
		log.Fatalln(err)
	}

	//
	// Create a providers changed event dispatcher
	//

	for {
		select {
		case update := <-listener:
			err = stream.Send(update)
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}

	return
}

func (broker *broker) dispatchProviders() {
	providers := broker.providerList()

	log.Printf("dispatchProviders: %d providers", len(providers.Items))

	broker.lmu.RLock()
	for _, ch := range broker.listener {
		ch <- providers
	}
	broker.lmu.RUnlock()
}

func (broker *broker) Register(stream pb.Broker_RegisterServer) (err error) {

	var ping *pb.Ping
	var id string

	for {
		ping, err = stream.Recv()
		if err != nil {
			break
		}

		switch data := ping.Cast.(type) {
		case *pb.Ping_Register:

			id = data.Register.ProviderId
			log.Printf("Register: connect id=%v", id)

			broker.pmu.Lock()
			broker.providers[id] = data.Register
			broker.pmu.Unlock()

			broker.dispatchProviders()

		case *pb.Ping_Util:

			if id == "" {
				continue
			}

			if data.Util == nil {
				continue
			}

			//log.Printf("Register: id=%v utilization=%v", id, data.Util.CpuUsage)

			broker.umu.Lock()
			broker.utilization[id] = data.Util
			broker.umu.Unlock()
		}
	}

	log.Printf("Register: disconnect id=%v", id)

	go func() {
		broker.pmu.Lock()
		delete(broker.providers, id)
		broker.pmu.Unlock()

		broker.dispatchProviders()
	}()

	go func() {
		broker.umu.Lock()
		delete(broker.utilization, id)
		broker.umu.Unlock()
	}()

	return
}
