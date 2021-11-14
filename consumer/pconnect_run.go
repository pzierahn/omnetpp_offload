package consumer

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"log"
	"time"
)

func (pConn *providerConnection) run(task *pb.SimulationRun) (ref *pb.StorageRef, err error) {
	runName := task.Config + "-" + task.RunNum
	log.Printf("[%s] %s start", pConn.id(), runName)

	start := time.Now()
	ref, err = pConn.provider.Run(pConn.ctx, task)
	if err != nil {
		log.Printf("[%s] %s error: %v", pConn.id(), runName, err)
	} else {
		log.Printf("[%s] %s finished (%v)", pConn.id(), runName, time.Now().Sub(start))
	}

	return
}
