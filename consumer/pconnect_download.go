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

	log.Printf("[%s] %s downloaded results %v in %v",
		pConn.id(), ref.Filename, simple.ByteSize(dlsize), time.Now().Sub(start))

	return
}

func (pConn *providerConnection) downloader() {
	// TODO:
}