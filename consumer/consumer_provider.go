package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
	"sync"
)

func (cons *consumer) availableProvider(broker pb.BrokerClient) {
	stream, err := broker.GetProviders(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalln(err)
	}

	for {
		providers, err := stream.Recv()
		if err != nil {
			log.Fatalf("exit providers updated event listener: %v", err)
		}

		log.Printf("providers updated event: %v", simple.PrettyString(providers.Items))

		var wg sync.WaitGroup
		connections := make(map[string]*connection)

		for _, prov := range providers.Items {

			cons.connMu.RLock()
			conn, ok := cons.connections[prov.ProviderId]
			cons.connMu.RUnlock()

			if ok {

				//
				// Connection already established
				//

				connections[prov.ProviderId] = conn

				continue
			}

			var mux sync.RWMutex

			wg.Add(1)
			go func(prov *pb.ProviderInfo) {
				defer wg.Done()

				conn, err := connect(prov)
				if err != nil {
					log.Println(err)
					return
				}

				err = conn.checkout(cons.simulation, cons.simulationTgz)
				if err != nil {
					log.Println(err)
					return
				}

				err = conn.setup(cons.simulation)
				if err != nil {
					log.Println(err)
					return
				}

				mux.Lock()
				connections[prov.ProviderId] = conn
				mux.Unlock()
			}(prov)
		}

		wg.Wait()

		cons.connMu.Lock()
		cons.connections = connections
		cons.connMu.Unlock()
	}
}
