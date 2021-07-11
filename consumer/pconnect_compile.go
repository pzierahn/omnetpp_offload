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

func (pConn *providerConnection) compileAndDownload(simulation *pb.Simulation) (err error) {

	arch := sysinfo.Signature(pConn.info.Arch)
	store := storage.FromClient(pConn.store)

	log.Printf("[%s] compileAndDownload: %s", pConn.name(), arch)

	var bin *pb.Binary
	bin, err = pConn.provider.Compile(context.Background(), simulation)
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

	log.Printf("[%s] compileAndDownload: %s done", pConn.name(), arch)

	return
}

func (pConn *providerConnection) setupExecutable(simulation *pb.Simulation) (err error) {

	arch := sysinfo.Signature(pConn.info.Arch)

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

	if !ok {
		err = pConn.compileAndDownload(simulation)
		return
	}

	binary := &checkoutObject{
		SimulationId: simulation.Id,
		Filename:     fmt.Sprintf("binary/%s.tgz", arch),
		Data:         buf,
	}

	err = pConn.checkout(binary)

	return
}
