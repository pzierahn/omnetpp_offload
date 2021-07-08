package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

func (conn *connection) run(task *pb.SimulationRun) (err error) {
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

	// TODO: Delete downloaded object from storage

	log.Printf("[%s] %s downloaded results %v in %v",
		conn.name(), runName, simple.ByteSize(uint64(buf.Len())), time.Now().Sub(endExec))

	dump := "/Users/patrick/Desktop/dump"
	err = ioutil.WriteFile(filepath.Join(dump, runName+".tgz"), buf.Bytes(), 0755)

	return
}
