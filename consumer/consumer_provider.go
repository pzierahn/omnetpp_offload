package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
	"sync"
)

func (cons *consumer) startConnector(broker pb.BrokerClient) {
	stream, err := broker.GetProviders(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalln(err)
	}

	var once sync.Once

	for {
		providers, err := stream.Recv()
		if err != nil {
			log.Fatalf("exit providers updated event listener: %v", err)
		}

		log.Printf("providers updated event: %v", simple.PrettyString(providers.Items))

		var wg sync.WaitGroup
		var mux sync.RWMutex
		connections := make(map[string]*connection)

		for _, prov := range providers.Items {

			cons.connMu.RLock()
			conn, ok := cons.connections[prov.ProviderId]
			cons.connMu.RUnlock()

			if ok {

				//
				// Connection already established, nothing to do
				//

				connections[prov.ProviderId] = conn

				continue
			}

			wg.Add(1)
			go func(prov *pb.ProviderInfo) {
				defer wg.Done()

				conn, err := connect(prov)
				if err != nil {
					log.Println(err)
					return
				}

				err = cons.init(conn)
				if err != nil {
					log.Println(err)
					return
				}

				once.Do(func() {
					log.Printf("[%s] list simulation run numbers", conn.name())

					runs, err := conn.provider.ListRunNums(context.Background(), &pb.Simulation{
						Id:        cons.simulation.Id,
						OppConfig: cons.simulation.OppConfig,

						// TODO: Fix 0 index
						Config: cons.config.SimulateConfigs[0],
					})
					if err != nil {
						log.Fatalln(err)
					}

					allocate := make([]*pb.SimulationRun, len(runs.Runs))
					for inx, run := range runs.Runs {
						allocate[inx] = &pb.SimulationRun{
							ConsumerId:   cons.consumerId,
							SimulationId: cons.simulation.Id,
							OppConfig:    cons.simulation.OppConfig,
							Config:       runs.Config,
							RunNum:       run,
						}
					}

					log.Printf("[%s] created %d jobs", conn.name(), len(allocate))

					cons.allocCond.L.Lock()
					// TODO: Remove slice cut!
					cons.allocate = allocate[:20]
					cons.allocCond.Broadcast()
					cons.allocCond.L.Unlock()
				})

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
