package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
	"time"
)

func (conn *connection) run(task *pb.Simulation) (err error) {
	runName := task.Config + "-" + task.RunNum
	log.Printf("[%s] %s start", conn.name(), runName)

	startExec := time.Now()

	resultRef, err := conn.provider.Run(context.Background(), task)
	if err != nil {
		return
	}

	endExec := time.Now()

	log.Printf("[%s] %s finished (%v)",
		conn.name(), runName, endExec.Sub(startExec))

	store := storage.FromClient(conn.store)
	buf, err := store.Download(resultRef)
	if err != nil {
		return
	}

	log.Printf("[%s] %s downloaded results %v in %v",
		conn.name(), runName, simple.ByteSize(uint64(buf.Len())), time.Now().Sub(endExec))

	return
}
