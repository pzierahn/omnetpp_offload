package consumer

import (
	"github.com/pzierahn/omnetpp_offload/eval"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"github.com/pzierahn/omnetpp_offload/storage"
	"log"
	"time"
)

func (pConn *providerConnection) download(ref *pb.StorageRef) (byt []byte, err error) {

	// Allow only one download process per provider.
	pConn.dmu.Lock()
	defer pConn.dmu.Unlock()

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

func (pConn *providerConnection) resultsDownloader(queue chan *download, sim *simulation) {
	for obj := range queue {
		buf, err := pConn.download(obj.ref)
		if err != nil {
			log.Printf("[%s] download failed: reschedule %+v", pConn.id(), obj.task)
			// Reschedule task.
			sim.queue.add(obj.task)
			return
		}

		done := eval.LogAction(eval.ActionExtract, obj.ref.Filename)
		err = simple.ExtractTarGz(sim.config.Path, buf)
		_ = done(err)
		if err != nil {
			log.Fatalf("cloudn't extract files: %v", err)
		}

		_, _ = pConn.store.Delete(pConn.ctx, obj.ref)

		sim.finished.Done()
	}
}
