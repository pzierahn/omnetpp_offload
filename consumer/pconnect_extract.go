package consumer

import (
	"github.com/pzierahn/omnetpp_offload/eval"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"github.com/pzierahn/omnetpp_offload/storage"
	"log"
	"time"
)

type fileMeta struct {
	SimulationId string
	Filename     string
	Data         []byte
}

func (pConn *providerConnection) extract(meta *fileMeta) (err error) {

	size := uint64(len(meta.Data))
	log.Printf("[%s] upload: %+v (%v)",
		pConn.id(), meta.Filename, simple.ByteSize(size))

	ui := make(chan storage.UploadProgress)
	defer close(ui)

	go func() {
		for info := range ui {
			log.Printf("[%s] upload: file=%s uploaded=%v percent=%0.2f%%",
				pConn.id(),
				meta.Filename,
				simple.ByteSize(info.Uploaded),
				info.Percent)
		}
	}()

	storeCli := storage.FromClient(pConn.store)

	start := time.Now()
	done := eval.Log(eval.Event{
		Activity: eval.ActivityUpload,
		Filename: meta.Filename,
	})

	upload := &storage.File{
		Bucket:   meta.SimulationId,
		Filename: meta.Filename,
		Data:     meta.Data,
	}
	ref, err := storeCli.Upload(upload, ui)
	done(err, size)

	if err != nil {
		return err
	}

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
