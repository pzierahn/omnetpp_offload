package consumer

import (
	"bytes"
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
	"time"
)

func (conn *connection) checkout(simulation *pb.Simulation, tgz []byte) (err error) {

	log.Printf("[%s] upload: %s (%v)",
		conn.name(), simulation.Id, simple.ByteSize(uint64(len(tgz))))

	startUpload := time.Now()

	ui := make(chan storage.UploadInfo)
	defer close(ui)

	go func() {
		for info := range ui {
			log.Printf("[%s] upload: simulation=%s uploaded=%v percent=%0.2f",
				conn.name(),
				simulation.Id,
				simple.ByteSize(info.Uploaded),
				float32(info.Uploaded)/float32(len(tgz)))
		}
	}()

	storeCli := storage.FromClient(conn.store)
	meta := storage.FileMeta{
		Bucket:   simulation.Id,
		Filename: "source.tar.gz",
	}

	ref, err := storeCli.Upload(bytes.NewReader(tgz), meta, ui)
	if err != nil {
		return
	}

	uploadTime := time.Now().Sub(startUpload)

	log.Printf("[%s] upload: %s finished (%v)", conn.name(), simulation.Id, uploadTime)

	log.Printf("[%s] checkout: %s...", conn.name(), simulation.Id)

	_, err = conn.provider.Checkout(context.Background(), &pb.Bundle{
		SimulationId: simulation.Id,
		Source:       ref,
	})

	log.Printf("[%s] checkout: %s done", conn.name(), simulation.Id)

	return
}
