package consumer

import (
	"bytes"
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
)

func (conn *connection) checkout(simulation *pb.Simulation, tgz []byte) (err error) {

	log.Printf("[%s] upload: %s (%v)",
		conn.name(), simulation.Id, simple.ByteSize(uint64(len(tgz))))

	storeCli := storage.FromClient(conn.store)
	var ref *pb.StorageRef
	ref, err = storeCli.Upload(bytes.NewReader(tgz), storage.FileMeta{
		Bucket:   simulation.Id,
		Filename: "source.tar.gz",
	})
	if err != nil {
		return
	}

	log.Printf("[%s] upload: %s done", conn.name(), simulation.Id)

	log.Printf("[%s] checkout: %s...", conn.name(), simulation.Id)

	_, err = conn.provider.Checkout(context.Background(), &pb.Bundle{
		SimulationId: simulation.Id,
		Source:       ref,
	})

	log.Printf("[%s] checkout: %s done", conn.name(), simulation.Id)

	return
}
