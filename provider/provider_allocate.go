package provider

import (
	"time"
)

func (prov *provider) register(simId string, allocRecv chan<- int) {

	newRecv := prov.newRecv
	newRecv.L.Lock()
	prov.allocRecvs[simId] = allocRecv
	newRecv.Broadcast()
	newRecv.L.Unlock()
}

func (prov *provider) unregister(simId string) {
	prov.mu.Lock()
	defer prov.mu.Unlock()

	delete(prov.allocRecvs, simId)
}

func (prov *provider) startAllocator(jobs int) {

	// Init the slot queue with the number of jobs.
	for inx := 0; inx < jobs; inx++ {
		prov.slots <- 1
	}

	newRecv := prov.newRecv

	for range prov.slots {

		newRecv.L.Lock()

		if len(prov.allocRecvs) == 0 {

			//
			// Wait for new allocation receivers.
			//

			newRecv.Wait()
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

		prov.allocRecvs[simId] <- 1

		newRecv.L.Unlock()
	}
}
