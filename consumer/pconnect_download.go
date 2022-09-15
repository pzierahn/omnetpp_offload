package consumer

import (
	"github.com/pzierahn/omnetpp_offload/eval"
	"github.com/pzierahn/omnetpp_offload/simple"
	"github.com/pzierahn/omnetpp_offload/storage"
	"log"
)

func (connect *providerConnection) download(dl *download) (byt []byte, err error) {

	// Allow only one download process per provider.
	connect.dmu.Lock()
	defer connect.dmu.Unlock()

	ref := dl.ref

	store := storage.FromClient(connect.store)

	done := eval.LogLocal(eval.Event{
		Activity:      eval.ActivityDownload,
		SimulationRun: dl.task,
		Filename:      ref.Filename,
	})

	byt, err = store.Download(connect.ctx, ref)
	size := uint64(len(byt))
	dur := done(err, size)

	if err != nil {
		log.Printf("[%s] error %v", connect.id(), err)
		return nil, err
	}

	log.Printf("[%s] download=%s size=%v time=%v", connect.id(), ref.Filename, simple.ByteSize(size), dur)

	return
}

func (connect *providerConnection) resultsDownloader(queue chan *download, sim *simulation) {
	for obj := range queue {
		buf, err := connect.download(obj)
		if err != nil {
			log.Printf("[%s] download failed: reschedule %+v", connect.id(), obj.task)
			// Reschedule task.
			sim.queue.add(obj.task)
			return
		}

		done := eval.Log(eval.Event{
			Activity:      eval.ActivityExtract,
			SimulationRun: obj.task,
			Filename:      obj.ref.Filename,
		})

		err = simple.ExtractTarGz(sim.config.Path, buf)
		done(err, 0)

		if err != nil {
			log.Fatalf("cloudn't extract files: %v", err)
		}

		_, _ = connect.store.Delete(connect.ctx, obj.ref)

		sim.finished.Done()
	}
}
