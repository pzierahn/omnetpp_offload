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
	var wg sync.WaitGroup
	var mux sync.RWMutex
	connections := make(map[string]*providerConnection)

	for {
		providers, err := stream.Recv()
		if err != nil {
			log.Fatalf("exit providers updated event listener: %v", err)
		}

		log.Printf("providers updated event: %v", simple.PrettyString(providers.Items))

		for _, prov := range providers.Items {

			cons.connMu.RLock()
			_, ok := connections[prov.ProviderId]
			cons.connMu.RUnlock()

			if ok {

				//
				// Connection already established, nothing to do
				//

				continue
			}

			wg.Add(1)
			go func(prov *pb.ProviderInfo) {
				defer wg.Done()

				var pconn *providerConnection
				pconn, err = cons.connect(prov)
				if err != nil {
					log.Println(prov.ProviderId, err)
					return
				}

				err = pconn.init(cons)
				if err != nil {
					log.Println(prov.ProviderId, err)
					return
				}

				mux.Lock()
				connections[prov.ProviderId] = pconn
				mux.Unlock()

				logProviderInfo(prov.ProviderId, prov)
			}(prov)
		}

		wg.Wait()

		for _, conn := range connections {
			once.Do(func() {
				log.Printf("[%s] list simulation run numbers", conn.name())

				allocate := make([]*pb.SimulationRun, 0)

				for _, conf := range cons.config.SimulateConfigs {
					runs, err := conn.provider.ListRunNums(context.Background(), &pb.Simulation{
						Id:        cons.simulation.Id,
						OppConfig: cons.simulation.OppConfig,
						Config:    conf,
					})
					if err != nil {
						log.Fatalln(err)
					}

					for _, run := range runs.Runs {
						allocate = append(allocate, &pb.SimulationRun{
							SimulationId: cons.simulation.Id,
							OppConfig:    cons.simulation.OppConfig,
							Config:       runs.Config,
							RunNum:       run,
						})
					}
				}

				log.Printf("[%s] created %d jobs", conn.name(), len(allocate))
				cons.allocate.add(allocate...)
				onInit <- cons.allocate.len()
			})

			break
		}
	}
}
