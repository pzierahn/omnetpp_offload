package consumer

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"sync"
)

func (sim *simulation) connect(prov *pb.ProviderInfo, once *sync.Once) {

	//
	// Phase 1: Connect to provider
	//

	cc, err := pconnect(sim.ctx, prov)
	if err != nil {
		log.Println(prov.ProviderId, err)
		return
	}

	pconn := &providerConnection{
		conn:     cc,
		info:     prov,
		provider: pb.NewProviderClient(cc),
		store:    pb.NewStorageClient(cc),
	}

	//
	// Phase 2: Deploy the simulation
	//

	err = pconn.deploy(sim)
	if err != nil {
		log.Println(prov.ProviderId, err)
		return
	}

	//
	// Phase 3: Execute the simulation
	//

	once.Do(func() {
		log.Printf("[%s] list simulation run numbers", pconn.id())

		tasks, err := pconn.collectTasks(sim)
		if err != nil {
			log.Fatalln(prov.ProviderId, err)
		}

		log.Printf("[%s] created %d jobs", pconn.id(), len(tasks))
		sim.queue.add(tasks...)
		sim.onInit <- sim.queue.len()
	})

	err = pconn.execute(sim)
	if err != nil {
		log.Println(prov.ProviderId, err)
		return
	}
}

func (sim *simulation) startConnector(bconn *grpc.ClientConn) {

	broker := pb.NewBrokerClient(bconn)
	stream, err := broker.Providers(sim.ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalln(err)
	}

	var once sync.Once
	var mux sync.RWMutex
	connections := make(map[string]bool)

	for {
		providers, err := stream.Recv()
		if err != nil {
			// TODO: Restart connector
			log.Printf("unable to recieve provider list update event: %v", err)
			return
		}

		log.Printf("providers update event: %v", simple.PrettyString(providers.Items))

		for _, prov := range providers.Items {

			mux.RLock()
			_, ok := connections[prov.ProviderId]
			mux.RUnlock()

			if ok {

				//
				// Connection already established, nothing to do
				//

				continue
			} else {
				//
				// Connect to provider
				//

				// TODO: Try to reconnect to the provider after fail

				mux.Lock()
				connections[prov.ProviderId] = true
				mux.Unlock()

				go sim.connect(prov, &once)
			}
		}
	}
}
