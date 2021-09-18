package consumer

import (
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

func (pConn *providerConnection) compileAndDownload(simulation *simulation) (err error) {

	arch := sysinfo.Signature(pConn.info.Arch)
	store := storage.FromClient(pConn.store)

	log.Printf("[%s] compile: %s", pConn.id(), arch)

	var bin *pb.Binary
	bin, err = pConn.provider.Compile(pConn.ctx, simulation.proto())
	if err != nil {
		return err
	}

	log.Printf("[%s] compile: %s done", pConn.id(), arch)

	start := time.Now()
	done := eval.LogTransfer(pConn.id(), eval.TransferDirectionDownload, bin.Ref.Filename)

	var byt []byte
	byt, err = store.Download(pConn.ctx, bin.Ref)
	if err != nil {
		return done(0, err)
	}

	size := uint64(len(byt))
	_ = done(size, err)

	log.Printf("[%s] compile: downloaded %s exe (%v in %v)",
		pConn.id(), arch, simple.ByteSize(size), time.Now().Sub(start))

	simulation.bmu.Lock()
	simulation.binaries[arch] = byt
	simulation.bmu.Unlock()

	return
}

func (pConn *providerConnection) setupExecutable(simulation *simulation) (err error) {

	arch := sysinfo.Signature(pConn.info.Arch)

	// TODO: Find an easy way to do this
	var lock *sync.Mutex
	simulation.amu.Lock()
	if aLock, ok := simulation.archLock[arch]; ok {
		lock = aLock
	} else {
		lock = &sync.Mutex{}
		simulation.archLock[arch] = lock
	}
	simulation.amu.Unlock()

	lock.Lock()
	defer lock.Unlock()

	simulation.bmu.RLock()
	buf, ok := simulation.binaries[arch]
	simulation.bmu.RUnlock()

	if !ok {
		err = pConn.compileAndDownload(simulation)
		return
	}

	binary := &checkoutObject{
		SimulationId: simulation.id,
		Filename:     fmt.Sprintf("binary/%s.tgz", arch),
		Data:         buf,
	}

	err = pConn.extract(binary)

	return
}
