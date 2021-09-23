package consumer

import (
	"github.com/pzierahn/project.go.omnetpp/eval"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"log"
	"time"
)

func (pConn *providerConnection) download(ref *pb.StorageRef) (byt []byte, err error) {

	start := time.Now()
	store := storage.FromClient(pConn.store)
	done := eval.LogTransfer(pConn.id(), eval.TransferDirectionDownload, ref.Filename)

	byt, err = store.Download(pConn.ctx, ref)
	if err != nil {
		log.Printf("[%s] error %v", pConn.id(), err)
		return byt, done(0, err)
	}

	dlsize := uint64(len(byt))
	_ = done(dlsize, nil)

	log.Printf("[%s] download=%s size=%v time=%v",
		pConn.id(), ref.Filename, simple.ByteSize(dlsize), time.Now().Sub(start))

	return
}

func (pConn *providerConnection) downloader(sim *simulation) {
	for obj := range pConn.downloadPipe {
		buf, err := pConn.download(obj.ref)
		if err != nil {
			log.Printf("[%s] download failed: reschedule %+v", pConn.id(), obj.task)
			// Add item back to taskQueue to send right allocation num
			sim.queue.add(obj.task)
			return
		}

		done := eval.LogAction(eval.ActionExtract, obj.ref.Filename)
		sim.extractResults(buf)
		_ = done(nil)

		_, _ = pConn.store.Delete(pConn.ctx, obj.ref)

		sim.finished.Done()
	}
}
