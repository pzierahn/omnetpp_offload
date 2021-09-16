package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"log"
	"time"
)

func (pConn *providerConnection) run(task *pb.SimulationRun) (ref *pb.StorageRef, err error) {
	runName := task.Config + "-" + task.RunNum
	log.Printf("[%s] %s start", pConn.id(), runName)

	start := time.Now()
	ref, err = pConn.provider.Run(pConn.ctx, task)
	if err != nil {
		log.Printf("[%s] error %v", pConn.id(), err)
	} else {
		log.Printf("[%s] %s finished (%v)", pConn.id(), runName, time.Now().Sub(start))
	}

	return
}
