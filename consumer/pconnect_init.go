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
				OppConfig:    cons.simulation.OppConfig,
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

	//
	// TODO: Set sessions attributes
	//
	session, err := pConn.provider.GetSession(cons.ctx, simulation)
	if err != nil {
		return
	}

	log.Printf("init: set execution deadline %s", session.Ttl.AsTime())

	source := &checkoutObject{
		SimulationId: simulation.Id,
		Filename:     "source.tgz",
		Data:         cons.simulationSource,
	}

	if err = pConn.checkout(source); err != nil {
		return
	}

	if err = pConn.setupExecutable(simulation); err != nil {
		return
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
