package consumer

import (
	"github.com/pzierahn/project.go.omnetpp/eval"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"io/ioutil"
	"log"
	"path/filepath"
)

func (pConn *providerConnection) run(task *pb.SimulationRun) (err error) {
	runName := task.Config + "-" + task.RunNum
	log.Printf("[%s] %s start", pConn.name(), runName)

	elog := eval.LogRun(eval.Run{
		Command:    eval.CommandExecution,
		ProviderId: pConn.name(),
		RunConfig:  task.Config,
		RunNumber:  task.RunNum,
	})

	resultRef, err := pConn.provider.Run(pConn.ctx, task)
	if err != nil {
		log.Printf("[%s] error %v", pConn.name(), err)
		return elog.Error(err)
	}

	exeDur := elog.Success()
	log.Printf("[%s] %s finished (%v)", pConn.name(), runName, exeDur)

	dlDur := eval.LogTransfer(eval.Transfer{
		ProviderId:       pConn.name(),
		Direction:        eval.TransferDirectionDownload,
		BytesTransferred: 0,
	})

	store := storage.FromClient(pConn.store)
	buf, err := store.Download(pConn.ctx, resultRef)
	if err != nil {
		log.Printf("[%s] error %v", pConn.name(), err)
		return dlDur.Error(err)
	}

	dlsize := uint64(buf.Len())
	transfer := dlDur.Success(dlsize)

	log.Printf("[%s] %s downloaded results %v in %v",
		pConn.name(), runName, simple.ByteSize(dlsize), transfer)

	// TODO: Extract files to the right place
	dump := "/Users/patrick/Desktop/dump"
	err = ioutil.WriteFile(filepath.Join(dump, runName+".tgz"), buf.Bytes(), 0755)
	if err != nil {
		log.Printf("[%s] error %v", pConn.name(), err)
		return
	}

	_, err = pConn.store.Delete(pConn.ctx, resultRef)

	return
}
