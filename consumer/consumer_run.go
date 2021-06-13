package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
	"sync"
	"time"
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

					runName := task.Config + "-" + task.RunNum
					log.Printf("[%s-%d] %s start", conn.info.ProviderId, agent, runName)

					startExec := time.Now()

					resultRef, err := conn.provider.Run(context.Background(), task)
					if err != nil {
						log.Fatalln(err)
					}

					endExec := time.Now()

					log.Printf("[%s-%d] %s finished (%v)",
						conn.info.ProviderId, agent, runName, endExec.Sub(startExec))

					store := storage.FromClient(conn.store)
					buf, err := store.Download(resultRef)
					if err != nil {
						log.Fatalln(err)
					}

					log.Printf("[%s-%d] %s downloaded results %d bytes (%v)",
						conn.info.ProviderId, agent, runName, buf.Len(), time.Now().Sub(endExec))

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
