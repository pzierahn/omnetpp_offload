package provider

import (
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
)

func (prov *provider) allocateSlots() {

	prov.mu.RLock()
	defer prov.mu.RUnlock()

	log.Printf("allocateSlots: slots=%d requesters=%d", prov.freeSlots, len(prov.requests))

	for cId, req := range prov.requests {

		if prov.freeSlots == 0 {
			return
		}

		var allocatedSlots uint32

		for _, num := range prov.assignments {
			allocatedSlots += num
		}

		assignable := prov.freeSlots - allocatedSlots
		log.Printf("allocateSlots: assignable=%d", assignable)

		if assignable == 0 {
			return
		}

		slots := simple.MathMinUint32(assignable, req)

		log.Printf("allocateSlots: assign cId=%s slots=%d freeSlots=%d", cId, slots, prov.freeSlots)

		prov.assignments[cId] = slots
		//prov.freeSlots -= slots
		prov.allocate[cId] <- slots
	}

	return
}
