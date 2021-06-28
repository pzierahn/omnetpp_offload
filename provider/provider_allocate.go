package provider

import (
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
	"sync/atomic"
)

func (prov *provider) allocator() {

	cond := prov.cond

	for {
		log.Printf("########## Lock")
		cond.L.Lock()
		log.Printf("########## Wait")
		cond.Wait()

		freeSlots := atomic.LoadInt32(&prov.freeSlots)
		log.Printf("allocator: freeSlots=%d", freeSlots)

		for cId, req := range prov.requests {

			if freeSlots == 0 {
				return
			}

			var allocatedSlots uint32

			for _, num := range prov.assignments {
				allocatedSlots += num
			}

			assignable := prov.slots - allocatedSlots
			log.Printf("allocator: assignable=%d allocatedSlots=%+v", assignable, prov.assignments)

			if assignable == 0 {
				return
			}

			// TODO: remove 1
			slots := simple.MathMinUint32(assignable, req, 1)

			log.Printf("allocator: assign cId=%s slots=%d freeSlots=%d", cId, slots, freeSlots)

			prov.assignments[cId] += slots
			//prov.freeSlots -= slots

			log.Printf("########## ch")
			prov.allocate[cId] <- slots
		}

		log.Printf("########## Unlock")
		cond.L.Unlock()
	}
}
