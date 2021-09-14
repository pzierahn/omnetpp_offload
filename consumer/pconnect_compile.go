package consumer

import (
	"bytes"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/eval"
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

func (pConn *providerConnection) compileAndDownload(simulation *pb.Simulation) (err error) {

	arch := sysinfo.Signature(pConn.info.Arch)
	store := storage.FromClient(pConn.store)

	log.Printf("[%s] compile: %s", pConn.id(), arch)

	var bin *pb.Binary
	bin, err = pConn.provider.Compile(pConn.ctx, simulation)
	if err != nil {
		//return ccDone.Error(err)
		return err
	}

	//duration := ccDone.Success()
	duration := "MISSING"
	log.Printf("[%s] compile: %s done (%v)", pConn.id(), arch, duration)

	start := time.Now()
	done := eval.LogTransfer(pConn.id(), eval.TransferDirectionDownload, bin.Ref.Filename)

	var buf bytes.Buffer
	buf, err = store.Download(pConn.ctx, bin.Ref)
	if err != nil {
		return done(0, err)
	}

	size := uint64(buf.Len())
	_ = done(size, err)

	log.Printf("[%s] compile: downloaded %s exe (%v in %v)",
		pConn.id(), arch, simple.ByteSize(size), time.Now().Sub(start))

	binaryMu.Lock()
	binaries[arch] = buf.Bytes()
	binaryMu.Unlock()

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

	err = pConn.extract(binary)

	return
}
