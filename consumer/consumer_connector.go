package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"sync"
)

func (sim *simulation) connect(prov *pb.ProviderInfo, once *sync.Once, onInit chan int32) {
	cc, err := pconnect(sim.ctx, prov)
	if err != nil {
		log.Println(prov.ProviderId, err)
		return
	}

	pconn := &providerConnection{
		conn:         cc,
		info:         prov,
		provider:     pb.NewProviderClient(cc),
		store:        pb.NewStorageClient(cc),
		downloadPipe: make(chan *download, 128),
	}

	err = pconn.init(sim)
	if err != nil {
		log.Println(prov.ProviderId, err)
		return
	}

	once.Do(func() {
		log.Printf("[%s] list simulation run numbers", pconn.id())

		tasks, err := pconn.collectTasks(sim)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("[%s] created %d jobs", pconn.id(), len(tasks))
		sim.queue.add(tasks...)
		onInit <- sim.queue.len()
	})
}

func (sim *simulation) startConnector(bconn *grpc.ClientConn, onInit chan int32) {

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
			log.Fatalf("exit providers updated event listener: %v", err)
		}

		log.Printf("providers updated event: %v", simple.PrettyString(providers.Items))

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

				go sim.connect(prov, &once, onInit)
			}
		}
	}
}
