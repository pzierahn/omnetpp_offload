package consumer

import (
	"fmt"
	"github.com/pzierahn/omnetpp_offload/eval"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"github.com/pzierahn/omnetpp_offload/storage"
	"github.com/pzierahn/omnetpp_offload/sysinfo"
	"log"
	"sync"
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

	done := eval.Log(eval.Event{
		Activity: eval.ActivityDownload,
		Filename: bin.Ref.Filename,
	})

	var byt []byte
	byt, err = store.Download(pConn.ctx, bin.Ref)

	size := uint64(len(byt))
	dur := done(err, size)

	if err != nil {
		return err
	}

	log.Printf("[%s] compile: downloaded %s exe (%v in %v)",
		pConn.id(), arch, simple.ByteSize(size), dur)

	simulation.bmu.Lock()
	simulation.binaries[arch] = byt
	simulation.bmu.Unlock()

	return
}

func (pConn *providerConnection) setupExecutable(simulation *simulation) (err error) {

	arch := sysinfo.Signature(pConn.info.Arch)

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

		//
		// Simulation executable is missing for providers architecture and OS.
		// Compile and download the executable.
		//

		err = pConn.compileAndDownload(simulation)
	} else {

		//
		// Simulation executable is already compiled for providers architecture and OS.
		//

		binary := &fileMeta{
			SimulationId: simulation.id,
			Filename:     fmt.Sprintf("binary/%s.tgz", arch),
			Data:         buf,
		}

		err = pConn.extract(binary)
	}

	return
}
