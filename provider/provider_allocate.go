package provider

import (
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
	"sync/atomic"
)

func (prov *provider) startAllocator() {

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

			var assignedSlots uint32

			for _, num := range prov.assignments {
				assignedSlots += num
			}

			assignable := prov.slots - assignedSlots
			log.Printf("allocator: assignable=%d assignments=%v", assignable, prov.assignments)

			if assignable == 0 {
				break
			}

			slots := simple.MathMinUint32(assignable, req)

			log.Printf("allocator: assign cId=%s slots=%d", cId, slots)

			prov.assignments[cId] += slots
			prov.allocate[cId] <- slots
		}

		cond.L.Unlock()
	}
}
