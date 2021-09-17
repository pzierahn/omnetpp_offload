package consumer

import (
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

func (pConn *providerConnection) extract(meta *checkoutObject) (err error) {

	size := uint64(len(meta.Data))
	log.Printf("[%s] upload: %+v (%v)",
		pConn.id(), meta.Filename, simple.ByteSize(size))

	ui := make(chan storage.UploadInfo)
	defer close(ui)

	go func() {
		for info := range ui {
			log.Printf("[%s] upload: file=%s uploaded=%v percent=%0.2f%%",
				pConn.id(),
				meta.Filename,
				simple.ByteSize(info.Uploaded),
				100*float32(info.Uploaded)/float32(len(meta.Data)))
		}
	}()

	storeCli := storage.FromClient(pConn.store)

	start := time.Now()
	done := eval.LogTransfer(pConn.id(), eval.TransferDirectionUpload, meta.Filename)

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

	log.Printf("[%s] upload: finished file=%s size=%s time=%v",
		pConn.id(), meta.Filename, simple.ByteSize(size), time.Now().Sub(start))

	log.Printf("[%s] extract: %s...", pConn.id(), meta.Filename)

	_, err = pConn.provider.Extract(pConn.ctx, &pb.Bundle{
		SimulationId: meta.SimulationId,
		Source:       ref,
	})

	log.Printf("[%s] extract: %s done",
		pConn.id(), meta.Filename)

	return
}
