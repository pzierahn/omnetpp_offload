package broker

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"log"
	"sync"
)

type broker struct {
	pb.UnimplementedBrokerServer
	mu        sync.RWMutex
	providers map[string]*pb.ProviderInfo
}

func (broker *broker) GetProviders(_ context.Context, _ *pb.Empty) (providers *pb.Providers, err error) {

	log.Printf("GetProviders:")

	providers = &pb.Providers{}

	broker.mu.RLock()
	for _, prov := range broker.providers {
		providers.Items = append(providers.Items, prov)
	}
	broker.mu.RUnlock()

	return
}

func (broker *broker) Register(_ context.Context, info *pb.ProviderInfo) (res *pb.Empty, err error) {

	log.Printf("Register: info=%v", info)

	broker.mu.Lock()
	broker.providers[info.ProviderId] = info
	broker.mu.Unlock()

	res = &pb.Empty{}
	return
}

func (broker *broker) Status(stream pb.Broker_StatusServer) (err error) {

	return
}
