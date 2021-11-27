package consumer

import (
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"log"
	"time"
)

func (connect *providerConnection) run(task *pb.SimulationRun) (ref *pb.StorageRef, err error) {
	runName := task.Config + "-" + task.RunNum
	log.Printf("[%s] %s start", connect.id(), runName)

	start := time.Now()
	ref, err = connect.provider.Run(connect.ctx, task)
	if err != nil {
		log.Printf("[%s] %s error: %v", connect.id(), runName, err)
	} else {
		log.Printf("[%s] %s finished (%v)", connect.id(), runName, time.Now().Sub(start))
	}

	return
}
