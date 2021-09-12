package consumer

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/eval"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
	"time"
)

type checkoutObject struct {
	SimulationId string
	Filename     string
	Data         []byte
}

func (pConn *providerConnection) checkout(meta *checkoutObject) (err error) {

	size := uint64(len(meta.Data))
	log.Printf("[%s] upload: %+v (%v)",
		pConn.name(), meta.Filename, simple.ByteSize(size))

	ui := make(chan storage.UploadInfo)
	defer close(ui)

	var info storage.UploadInfo

	go func() {
		for info = range ui {
			log.Printf("[%s] upload: file=%s uploaded=%v percent=%0.2f%%",
				pConn.name(),
				meta.Filename,
				simple.ByteSize(info.Uploaded),
				100*float32(info.Uploaded)/float32(len(meta.Data)))
		}
	}()

	storeCli := storage.FromClient(pConn.store)

	start := time.Now()
	done := eval.LogTransfer(pConn.name(), eval.TransferDirectionUpload, meta.Filename)

	upload := &storage.FileMeta{
		Bucket:   meta.SimulationId,
		Filename: meta.Filename,
		Data:     meta.Data,
	}
	ref, err := storeCli.Upload(upload, ui)
	if err != nil {
		return done(0, err)
	}

	_ = done(size, nil)

	log.Printf("[%s] upload: finished file=%s packages=%d size=%s time=%v",
		pConn.name(), meta.Filename, info.Parcels, simple.ByteSize(size), time.Now().Sub(start))

	//checkoutDuration := eval.LogRun(eval.Run{
	//	Command:    eval.ActionCheckout,
	//	ProviderId: pConn.name(),
	//})

	log.Printf("[%s] checkout: %s...", pConn.name(), meta.Filename)

	_, err = pConn.provider.Extract(context.Background(), &pb.Bundle{
		SimulationId: meta.SimulationId,
		Source:       ref,
	})

	//checkoutDur := checkoutDuration.Success()
	checkoutDur := "MISSING"

	log.Printf("[%s] checkout: %s done (%v)",
		pConn.name(), meta.Filename, checkoutDur)

	//
	// TODO: Delete checked-out refs
	//

	return
}
