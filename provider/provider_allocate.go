package provider

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
)

func (prov *provider) dropSession(id simulationId) {

	log.Printf("dropSession: simulationId=%v", id)

	delete(prov.allocate, id)
	delete(prov.requests, id)
	delete(prov.assignments, id)
	delete(prov.sessions, id)

	// Clean up and remove simulation (delete simulation bucket)
	_, _ = prov.store.Drop(nil, &pb.BucketRef{Bucket: id})

	dir := filepath.Join(cachePath, id)
	_ = os.RemoveAll(dir)
}

func (prov *provider) allocator() {

	cond := prov.cond

	for {
		cond.L.Lock()
		cond.Wait()

		freeSlots := atomic.LoadInt32(&prov.freeSlots)
		log.Printf("allocator: freeSlots=%d requests=%v", freeSlots, prov.requests)

		for cId, req := range prov.requests {

			if freeSlots == 0 {
				break
			}

			var allocatedSlots uint32

			for _, num := range prov.assignments {
				allocatedSlots += num
			}

			assignable := prov.slots - allocatedSlots
			log.Printf("allocator: assignable=%d allocatedSlots=%v", assignable, prov.assignments)

			if assignable == 0 {
				break
			}

			// TODO: remove 1
			slots := simple.MathMinUint32(assignable, req, 1)

			log.Printf("allocator: assign cId=%s slots=%d", cId, slots)

			prov.assignments[cId] += slots
			prov.allocate[cId] <- slots
		}

		cond.L.Unlock()
	}
}
