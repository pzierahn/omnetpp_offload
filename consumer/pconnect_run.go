package consumer

import (
	"bytes"
	"github.com/pzierahn/project.go.omnetpp/eval"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
	"time"
)

func (pConn *providerConnection) run(task *pb.SimulationRun, config *Config) (err error) {
	runName := task.Config + "-" + task.RunNum
	log.Printf("[%s] %s start", pConn.id(), runName)

	start := time.Now()
	resultRef, err := pConn.provider.Run(pConn.ctx, task)
	if err != nil {
		log.Printf("[%s] error %v", pConn.id(), err)
		return err
	}

	log.Printf("[%s] %s finished (%v)", pConn.id(), runName, time.Now().Sub(start))

	done := eval.LogTransfer(pConn.id(), eval.TransferDirectionDownload, resultRef.Filename)
	store := storage.FromClient(pConn.store)

	start = time.Now()
	buf, err := store.Download(pConn.ctx, resultRef)
	if err != nil {
		log.Printf("[%s] error %v", pConn.id(), err)
		return done(0, err)
	}

	dlsize := uint64(buf.Len())
	_ = done(dlsize, nil)

	log.Printf("[%s] %s downloaded results %v in %v",
		pConn.id(), runName, simple.ByteSize(dlsize), time.Now().Sub(start))

	//
	// Extract files to the right place
	//

	err = simple.ExtractTarGz(config.Path, bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Printf("[%s] error %v", pConn.id(), err)
		return
	}

	_, err = pConn.store.Delete(pConn.ctx, resultRef)

	return
}
