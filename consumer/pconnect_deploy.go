package consumer

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"log"
)

func (connect *providerConnection) collectTasks(sim *simulation) (tasks []*pb.SimulationRun, err error) {

	for _, conf := range sim.config.SimulateConfigs {

		var runs *pb.SimulationRunList
		runs, err = connect.provider.ListRunNums(connect.ctx, &pb.Simulation{
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

func (connect *providerConnection) deploy(sim *simulation) (err error) {

	ctx := sim.ctx
	connect.ctx = ctx

	session, err := connect.provider.GetSession(ctx, sim.proto())
	if err != nil {
		return
	}

	log.Printf("deploy: deadline=%s source=%v exec=%v",
		session.Ttl.AsTime(), session.SourceExtracted, session.ExecutableExtracted)

	if !session.SourceExtracted {
		source := &fileMeta{
			SimulationId: sim.id,
			Filename:     "source.tgz",
			Data:         sim.source,
		}

		if err = connect.extract(source); err != nil {
			return
		}

		session.SourceExtracted = true
		session, _ = connect.provider.SetSession(ctx, session)
	}

	if !session.ExecutableExtracted {
		if err = connect.setupExecutable(sim); err != nil {
			return
		}

		session.ExecutableExtracted = true
		session, _ = connect.provider.SetSession(ctx, session)
	}

	return
}
