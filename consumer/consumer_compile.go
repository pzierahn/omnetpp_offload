package consumer

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"github.com/pzierahn/project.go.omnetpp/sysinfo"
	"log"
	"sync"
	"time"
)

var archMu sync.Mutex
var archLock = make(map[string]*sync.Mutex)

var binaryMu sync.RWMutex
var binaries = make(map[string][]byte)

func (conn *connection) compile(simulation *pb.Simulation) (err error) {

	arch := sysinfo.Signature(conn.info.Arch)
	store := storage.FromClient(conn.store)

	log.Printf("[%s] compile: %s", conn.name(), arch)

	var bin *pb.Binary
	bin, err = conn.provider.Compile(context.Background(), simulation)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	buf, err = store.Download(bin.Ref)
	if err != nil {
		return
	}

	binaryMu.Lock()
	binaries[arch] = buf.Bytes()
	binaryMu.Unlock()

	log.Printf("[%s] compile: %s done", conn.name(), arch)

	return
}

func (conn *connection) uploadBinary(simulation *pb.Simulation, buf []byte) (err error) {

	log.Printf("[%s] uploadBinary:", conn.name())

	arch := sysinfo.Signature(conn.info.Arch)
	store := storage.FromClient(conn.store)

	meta := storage.FileMeta{
		Bucket:   simulation.Id,
		Filename: fmt.Sprintf("binary/%s.tgz", arch),
	}

	startUpload := time.Now()

	ui := make(chan storage.UploadInfo)
	defer close(ui)

	go func() {
		for info := range ui {
			log.Printf("[%s] upload: simulation=%s uploaded=%v percent=%0.2f",
				conn.name(),
				simulation.Id,
				simple.ByteSize(info.Uploaded),
				float32(info.Uploaded)/float32(len(buf)))
		}
	}()

	ref, err := store.Upload(bytes.NewReader(buf), meta, ui)
	if err != nil {
		return
	}

	_, err = conn.provider.Checkout(context.Background(), &pb.Bundle{
		SimulationId: simulation.Id,
		Source:       ref,
	})

	dur := time.Now().Sub(startUpload)
	log.Printf("[%s] uploadBinary: finished in %v", conn.name(), dur)

	return
}

func (conn *connection) setup(simulation *pb.Simulation) (err error) {

	arch := sysinfo.Signature(conn.info.Arch)

	// TODO: Find an easy way to do this
	var lock *sync.Mutex
	archMu.Lock()
	if aLock, ok := archLock[arch]; ok {
		lock = aLock
	} else {
		lock = &sync.Mutex{}
		archLock[arch] = lock
	}
	archMu.Unlock()

	lock.Lock()
	defer lock.Unlock()

	binaryMu.RLock()
	buf, ok := binaries[arch]
	binaryMu.RUnlock()

	if ok {
		err = conn.uploadBinary(simulation, buf)
	} else {
		err = conn.compile(simulation)
	}

	return
}
