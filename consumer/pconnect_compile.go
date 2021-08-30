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
)

var archMu sync.Mutex
var archLock = make(map[string]*sync.Mutex)

var binaryMu sync.RWMutex
var binaries = make(map[string][]byte)

func (pConn *providerConnection) compileAndDownload(simulation *pb.Simulation) (err error) {

	arch := sysinfo.Signature(pConn.info.Arch)
	store := storage.FromClient(pConn.store)

	log.Printf("[%s] compile: %s", pConn.name(), arch)

	ccDone := eval.LogRun(eval.Run{
		Command:    eval.CommandCompile,
		ProviderId: pConn.name(),
	})

	var bin *pb.Binary
	bin, err = pConn.provider.Compile(pConn.ctx, simulation)
	if err != nil {
		return ccDone.Error(err)
	}

	duration := ccDone.Success()
	log.Printf("[%s] compile: %s done (%v)", pConn.name(), arch, duration)

	//downDone := eval.RecordDuration(eval.DurationDownload, pConn.name())
	downDone := eval.LogTransfer(eval.Transfer{
		ProviderId: pConn.name(),
		Direction:  eval.TransferDirectionDownload,
	})

	var buf bytes.Buffer
	buf, err = store.Download(pConn.ctx, bin.Ref)
	if err != nil {
		return
	}

	size := uint64(buf.Len())
	dlDuration := downDone.Success(size)

	log.Printf("[%s] compile: downloaded %s exe (%v in %v)",
		pConn.name(), arch, simple.ByteSize(size), dlDuration)

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

	err = pConn.checkout(binary)

	return
}
