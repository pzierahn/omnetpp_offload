package provider

import "sort"

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
			cond.Wait()
		}

		simIds := make([]string, len(prov.allocRecvs))
		var inx int
		for simId := range prov.allocRecvs {
			simIds[inx] = simId
			inx++
		}

		sort.SliceStable(simIds, func(i, j int) bool {
			simId1 := simIds[i]
			simId2 := simIds[j]

			duration1 := prov.executionTimes[simId1]
			duration2 := prov.executionTimes[simId2]

			return duration1 < duration2
		})

		simId := simIds[0]
		ch := prov.allocRecvs[simId]
		ch <- 1

		cond.L.Unlock()
	}
}
