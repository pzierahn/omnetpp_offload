package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
)

func (cons *consumer) zipSource() (err error) {

	log.Printf("zipping simulation source: %s", cons.config.Path)

	cons.simulationTgz, err = simple.TarGz(cons.config.Path, cons.simulation.Id, cons.config.Ignore...)

	return
}

func (cons *consumer) checkoutSimulations() (err error) {

	cons.connMu.RLock()
	defer cons.connMu.RUnlock()

	for id, conn := range cons.connections {
		log.Printf("upload: %s to %s (%d bytes)", cons.simulation.Id, id, cons.simulationTgz.Len())

		storeCli := storage.FromClient(conn.store)
		var ref *pb.StorageRef
		ref, err = storeCli.Upload(&cons.simulationTgz, storage.FileMeta{
			Bucket:   cons.simulation.Id,
			Filename: "source.tar.gz",
		})
		if err != nil {
			return
		}

		log.Printf("checkout: %s on %s", cons.simulation.Id, id)

		_, err = conn.provider.Checkout(context.Background(), &pb.Bundle{
			SimulationId: cons.simulation.Id,
			Source:       ref,
		})
		if err != nil {
			return
		}
	}

	return
}
