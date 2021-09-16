package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"log"
)

func (pConn *providerConnection) collectTasks(cons *consumer) (tasks []*pb.SimulationRun, err error) {

	for _, conf := range cons.config.SimulateConfigs {

		var runs *pb.SimulationRuns
		runs, err = pConn.provider.ListRunNums(pConn.ctx, &pb.Simulation{
			Id:        cons.simulation.Id,
			OppConfig: cons.simulation.OppConfig,
			Config:    conf,
		})
		if err != nil {
			return
		}

		for _, run := range runs.Runs {
			tasks = append(tasks, &pb.SimulationRun{
				SimulationId: cons.simulation.Id,
				Config:       runs.Config,
				RunNum:       run,
			})
		}
	}

	return
}

func (pConn *providerConnection) init(cons *consumer) (err error) {

	simulation := cons.simulation

	pConn.ctx = cons.ctx

	session, err := pConn.provider.GetSession(cons.ctx, simulation)
	if err != nil {
		return
	}

	log.Printf("init: deadline=%s source=%v exec=%v",
		session.Ttl.AsTime(), session.SourceExtracted, session.ExecutableExtracted)

	source := &checkoutObject{
		SimulationId: simulation.Id,
		Filename:     "source.tgz",
		Data:         cons.simulationSource,
	}

	if !session.SourceExtracted {
		if err = pConn.extract(source); err != nil {
			return
		}

		session.SourceExtracted = true
		session, _ = pConn.provider.SetSession(cons.ctx, session)
	}

	if !session.ExecutableExtracted {
		if err = pConn.setupExecutable(simulation); err != nil {
			return
		}

		session.ExecutableExtracted = true
		session, _ = pConn.provider.SetSession(cons.ctx, session)
	}

	stream, err := pConn.provider.Allocate(cons.ctx)
	if err != nil {
		return
	}

	go pConn.allocationHandler(stream, cons)

	err = pConn.sendAllocationRequest(stream, cons)
	if err != nil {
		return
	}

	go cons.allocate.onUpdate(func() (cancel bool) {
		err = pConn.sendAllocationRequest(stream, cons)
		if err != nil {
			log.Println(err)
			return true
		}

		return false
	})

	return
}
