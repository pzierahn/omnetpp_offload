package consumer

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"github.com/pzierahn/project.go.omnetpp/sysinfo"
	"log"
)

func (cons *consumer) compile() (err error) {

	cons.connMu.RLock()

	bins := make(map[string][]byte)

	for id, conn := range cons.connections {

		archSig := sysinfo.Signature(conn.info.Arch)
		storeCli := storage.FromClient(conn.store)

		if buf, ok := bins[archSig]; ok {
			log.Printf("compile: id=%s arch=%s cached", id, archSig)

			var ref *pb.StorageRef
			ref, err = storeCli.Upload(bytes.NewReader(buf), storage.FileMeta{
				Bucket:   cons.simulation.Id,
				Filename: fmt.Sprintf("binary/%s.tgz", archSig),
			})
			if err != nil {
				return
			}

			_, err = conn.provider.Checkout(context.Background(), &pb.Bundle{
				SimulationId: cons.simulation.Id,
				Source:       ref,
			})
			if err != nil {
				return
			}

			continue
		}

		log.Printf("compile: id=%s arch=%s", id, archSig)
		var bin *pb.Binary
		bin, err = conn.provider.Compile(context.Background(), cons.simulation)
		if err != nil {
			return
		}

		var buf bytes.Buffer
		buf, err = storeCli.Download(bin.Ref)
		if err != nil {
			return
		}

		bins[archSig] = buf.Bytes()
	}

	defer cons.connMu.RUnlock()

	return
}
