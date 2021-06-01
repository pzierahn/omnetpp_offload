package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
	"sync"
)

func (cons *consumer) run() (err error) {

	cons.connMu.RLock()
	defer cons.connMu.RUnlock()

	var runs *pb.SimulationRuns

	for _, conn := range cons.connections {
		runs, err = conn.provider.ListRunNums(context.Background(), &pb.Simulation{
			Id:        cons.simulation.Id,
			OppConfig: cons.simulation.OppConfig,

			// TODO: Fix 0 index this
			Config: cons.config.SimulateConfigs[0],
		})
		if err != nil {
			return
		}

		break
	}

	var wg sync.WaitGroup
	tasks := make(chan *pb.Simulation)
	defer close(tasks)

	for _, conn := range cons.connections {
		for inx := uint32(0); inx < conn.info.NumCPUs; inx++ {
			go func(conn *connection, agent uint32) {
				for task := range tasks {
					log.Printf("[%s-%d] %v-%v", conn.info.ProviderId, agent, task.Config, task.RunNum)
					resultRef, err := conn.provider.Run(context.Background(), task)
					if err != nil {
						log.Fatalln(err)
					}

					log.Printf("[%s-%d] result: %v", conn.info.ProviderId, agent, resultRef)

					store := storage.FromClient(conn.store)
					buf, err := store.Download(resultRef)
					if err != nil {
						log.Fatalln(err)
					}

					log.Printf("[%s-%d] downloaded %d bytes", conn.info.ProviderId, agent, buf.Len())

					wg.Done()
				}
			}(conn, inx)
		}
	}

	for _, num := range runs.Runs {
		wg.Add(1)
		tasks <- &pb.Simulation{
			Id:        cons.simulation.Id,
			OppConfig: cons.simulation.OppConfig,
			Config:    runs.Config,
			RunNum:    num,
		}
	}

	wg.Wait()

	return
}
