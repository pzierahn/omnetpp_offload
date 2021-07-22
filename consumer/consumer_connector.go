package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
	"sync"
)

func (cons *consumer) startConnector(onInit chan int32) {

	broker := pb.NewBrokerClient(cons.bconn)
	stream, err := broker.GetProviders(context.Background(), &pb.Empty{})
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
			}

			//
			// TODO: Try to reconnect to the provider after fail
			//

			mux.Lock()
			connections[prov.ProviderId] = true
			mux.Unlock()

			go func(prov *pb.ProviderInfo) {

				var pconn *providerConnection
				pconn, err = cons.connect(prov)
				if err != nil {
					log.Println(prov.ProviderId, err)
					return
				}

				stat.SetInfo(prov.ProviderId, prov)

				err = pconn.init(cons)
				if err != nil {
					log.Println(prov.ProviderId, err)
					return
				}

				once.Do(func() {
					log.Printf("[%s] list simulation run numbers", pconn.name())

					tasks, err := pconn.collectTasks(cons)
					if err != nil {
						log.Fatalln(err)
					}

					log.Printf("[%s] created %d jobs", pconn.name(), len(tasks))
					cons.allocate.add(tasks...)
					onInit <- cons.allocate.len()
				})
			}(prov)
		}
	}
}
