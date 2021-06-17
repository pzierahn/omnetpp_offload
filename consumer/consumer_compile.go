package consumer

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"github.com/pzierahn/project.go.omnetpp/sysinfo"
	"log"
	"sync"
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

	var ref *pb.StorageRef
	ref, err = store.Upload(bytes.NewReader(buf), storage.FileMeta{
		Bucket:   simulation.Id,
		Filename: fmt.Sprintf("binary/%s.tgz", arch),
	})
	if err != nil {
		return
	}

	_, err = conn.provider.Checkout(context.Background(), &pb.Bundle{
		SimulationId: simulation.Id,
		Source:       ref,
	})

	log.Printf("[%s] uploadBinary: done", conn.name())

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
