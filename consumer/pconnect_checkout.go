package consumer

import (
	"context"
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

	log.Printf("[%s] upload: %+v (%v)",
		pConn.name(), meta.Filename, simple.ByteSize(uint64(len(meta.Data))))

	startUpload := time.Now()

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

	upload := &storage.FileMeta{
		Bucket:   meta.SimulationId,
		Filename: meta.Filename,
		Data:     meta.Data,
	}
	ref, err := storeCli.Upload(upload, ui)
	if err != nil {
		return
	}

	uploadTime := time.Now().Sub(startUpload)
	startCheckout := time.Now()

	log.Printf("[%s] upload: finished file=%s packages=%d time=%v",
		pConn.name(), meta.Filename, info.Parcels, uploadTime)

	log.Printf("[%s] checkout: %s...", pConn.name(), meta.Filename)

	_, err = pConn.provider.Checkout(context.Background(), &pb.Bundle{
		SimulationId: meta.SimulationId,
		Source:       ref,
	})

	log.Printf("[%s] checkout: %s done (%v)",
		pConn.name(), meta.Filename, time.Now().Sub(startCheckout))

	//
	// TODO: Delete checked-out refs
	//

	return
}
