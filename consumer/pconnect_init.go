package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"log"
)

func (pConn *providerConnection) collectTasks(sim *simulation) (tasks []*pb.SimulationRun, err error) {

	for _, conf := range sim.config.SimulateConfigs {

		var runs *pb.SimulationRunList
		runs, err = pConn.provider.ListRunNums(pConn.ctx, &pb.Simulation{
			Id:        sim.id,
			OppConfig: sim.config.OppConfig,
			Config:    conf,
		})
		if err != nil {
			return
		}

		tasks = append(tasks, runs.Items...)
	}

	return
}

func (pConn *providerConnection) init(sim *simulation) (err error) {

	ctx := sim.ctx
	pConn.ctx = ctx

	session, err := pConn.provider.GetSession(ctx, sim.proto())
	if err != nil {
		return
	}

	log.Printf("init: deadline=%s source=%v exec=%v",
		session.Ttl.AsTime(), session.SourceExtracted, session.ExecutableExtracted)

	source := &checkoutObject{
		SimulationId: sim.id,
		Filename:     "source.tgz",
		Data:         sim.source,
	}

	if !session.SourceExtracted {
		if err = pConn.extract(source); err != nil {
			return
		}

		session.SourceExtracted = true
		session, _ = pConn.provider.SetSession(ctx, session)
	}

	if !session.ExecutableExtracted {
		if err = pConn.setupExecutable(sim); err != nil {
			return
		}

		session.ExecutableExtracted = true
		session, _ = pConn.provider.SetSession(ctx, session)
	}

	go pConn.downloader(1, sim)

	stream, err := pConn.provider.Allocate(ctx)
	if err != nil {
		return
	}

	go pConn.allocationHandler(stream, sim)

	go sim.queue.onChange(func() (cancel bool) {
		err = pConn.sendAllocationRequest(stream, sim)
		if err != nil {
			log.Println(err)
			return true
		}

		return false
	})

	return
}
