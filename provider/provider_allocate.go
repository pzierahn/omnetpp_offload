package provider

import (
	"time"
)

func (prov *provider) register(simId string, allocRecv chan<- int) {

	cond := prov.cond
	cond.L.Lock()
	prov.allocRecvs[simId] = allocRecv
	cond.Broadcast()
	cond.L.Unlock()
}

func (prov *provider) unregister(simId string) {
	prov.mu.Lock()
	defer prov.mu.Unlock()

	delete(prov.allocRecvs, simId)
}

func (prov *provider) startAllocator() {

	cond := prov.cond

	for range prov.slots {

		cond.L.Lock()

		if len(prov.allocRecvs) == 0 {

			//
			// Wait for new allocation receivers.
			//

			cond.Wait()
		}

		var simId string
		var lowest = time.Duration(-1)

		for id := range prov.allocRecvs {
			duration := prov.executionTimes[id]
			if lowest <= 0 || lowest > duration {
				simId = id
				lowest = duration
			}
		}

		ch := prov.allocRecvs[simId]
		ch <- 1

		cond.L.Unlock()
	}
}
