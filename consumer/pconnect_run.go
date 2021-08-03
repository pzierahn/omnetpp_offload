package consumer

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/statistic"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

func (pConn *providerConnection) run(task *pb.SimulationRun) (err error) {
	runName := task.Config + "-" + task.RunNum
	log.Printf("[%s] %s start", pConn.name(), runName)

	startExec := time.Now()

	resultRef, err := pConn.provider.Run(pConn.ctx, task)
	if err != nil {
		return
	}

	endExec := time.Now()

	log.Printf("[%s] %s finished (%v)",
		pConn.name(), runName, endExec.Sub(startExec))

	stat.GetExecution(pConn.name()).Add(endExec.Sub(startExec))

	store := storage.FromClient(pConn.store)
	buf, err := store.Download(pConn.ctx, resultRef)
	if err != nil {
		return
	}

	transfer := time.Now().Sub(endExec)
	size := uint64(buf.Len())

	log.Printf("[%s] %s downloaded results %v in %v",
		pConn.name(), runName, simple.ByteSize(size), transfer)

	stat.GetDownload(pConn.name()).Add(statistic.Transfer{
		Duration: transfer,
		Size:     size,
	})

	// TODO: Extract files to the right place
	dump := "/Users/patrick/Desktop/dump"
	err = ioutil.WriteFile(filepath.Join(dump, runName+".tgz"), buf.Bytes(), 0755)
	if err != nil {
		return
	}

	_, err = pConn.store.Delete(pConn.ctx, resultRef)

	return
}
