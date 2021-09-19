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

func (pConn *providerConnection) deploy(sim *simulation) (err error) {

	ctx := sim.ctx
	pConn.ctx = ctx

	session, err := pConn.provider.GetSession(ctx, sim.proto())
	if err != nil {
		return
	}

	log.Printf("deploy: deadline=%s source=%v exec=%v",
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

	return
}
