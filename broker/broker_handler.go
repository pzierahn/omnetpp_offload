package broker

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"log"
	"sync"
)

type providerId = string

type broker struct {
	pb.UnimplementedBrokerServer
	pmu         sync.RWMutex
	providers   map[providerId]*pb.ProviderInfo
	umu         sync.RWMutex
	utilization map[providerId]*pb.Utilization
}

func (broker *broker) GetProviders(_ context.Context, _ *pb.Empty) (providers *pb.Providers, err error) {

	log.Printf("GetProviders:")

	providers = &pb.Providers{}

	broker.pmu.RLock()
	for _, prov := range broker.providers {
		providers.Items = append(providers.Items, prov)
	}
	broker.pmu.RUnlock()

	return
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
	}()

	go func() {
		broker.umu.Lock()
		delete(broker.utilization, id)
		broker.umu.Unlock()
	}()

	return
}
