package consumer

import (
	"github.com/pzierahn/omnetpp_offload/eval"
	"github.com/pzierahn/omnetpp_offload/simple"
	"github.com/pzierahn/omnetpp_offload/storage"
	"log"
	"time"
)

func (pConn *providerConnection) download(dl *download) (byt []byte, err error) {

	// Allow only one download process per provider.
	pConn.dmu.Lock()
	defer pConn.dmu.Unlock()

	ref := dl.ref

	start := time.Now()
	store := storage.FromClient(pConn.store)
	//done := eval.LogTransfer(pConn.id(), eval.TransferDirectionDownload, ref.Filename)
	done := eval.Log(eval.Event{
		DeviceId:      pConn.id(),
		Activity:      eval.ActivityDownload,
		SimulationRun: dl.task,
		Filename:      ref.Filename,
	})

	byt, err = store.Download(pConn.ctx, ref)
	size := uint64(len(byt))
	done(err, size)

	if err != nil {
		log.Printf("[%s] error %v", pConn.id(), err)
		return nil, err
	}

	log.Printf("[%s] download=%s size=%v time=%v",
		pConn.id(), ref.Filename, simple.ByteSize(size), time.Now().Sub(start))

	return
}

func (pConn *providerConnection) resultsDownloader(queue chan *download, sim *simulation) {
	for obj := range queue {
		buf, err := pConn.download(obj)
		if err != nil {
			log.Printf("[%s] download failed: reschedule %+v", pConn.id(), obj.task)
			// Reschedule task.
			sim.queue.add(obj.task)
			return
		}

		done := eval.Log(eval.Event{
			DeviceId:      pConn.id(),
			Activity:      eval.ActivityExtract,
			SimulationRun: obj.task,
			Filename:      obj.ref.Filename,
		})

		err = simple.ExtractTarGz(sim.config.Path, buf)
		done(err, 0)

		if err != nil {
			log.Fatalf("cloudn't extract files: %v", err)
		}

		_, _ = pConn.store.Delete(pConn.ctx, obj.ref)

		sim.finished.Done()
	}
}
