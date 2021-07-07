package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
	"time"
)

func (conn *connection) checkout(meta storage.FileMeta) (err error) {

	simulation := conn.simulation

	log.Printf("[%s] upload: %s (%v)",
		conn.name(), simulation.Id, simple.ByteSize(uint64(len(meta.Data))))

	startUpload := time.Now()

	ui := make(chan storage.UploadInfo)
	defer close(ui)

	var info storage.UploadInfo

	go func() {
		for info = range ui {
			log.Printf("[%s] upload: simulation=%s uploaded=%v percent=%0.2f%%",
				conn.name(),
				simulation.Id,
				simple.ByteSize(info.Uploaded),
				100*float32(info.Uploaded)/float32(len(meta.Data)))
		}
	}()

	storeCli := storage.FromClient(conn.store)

	ref, err := storeCli.Upload(meta, ui)
	if err != nil {
		return
	}

	uploadTime := time.Now().Sub(startUpload)

	log.Printf("[%s] upload: %s finished packages=%d time=%v", conn.name(), simulation.Id, info.Parcels, uploadTime)

	log.Printf("[%s] checkout: %s...", conn.name(), simulation.Id)

	_, err = conn.provider.Checkout(context.Background(), &pb.Bundle{
		SimulationId: simulation.Id,
		Source:       ref,
	})

	log.Printf("[%s] checkout: %s done", conn.name(), simulation.Id)

	return
}
