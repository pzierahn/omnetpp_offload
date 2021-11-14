package provider

import "github.com/pzierahn/omnetpp_offload/simple"

func startWatchers(prov *provider) {
	simple.Watch("/sessions", func() interface{} {
		prov.mu.RLock()
		defer prov.mu.RUnlock()

		return prov.sessions
	})
	simple.Watch("/executionTimes", func() interface{} {
		prov.mu.RLock()
		defer prov.mu.RUnlock()

		data := make(map[string]string)

		for id, dur := range prov.executionTimes {
			data[id] = dur.String()
		}

		return data
	})
	simple.Watch("/allocRecvs", func() interface{} {
		prov.mu.RLock()
		defer prov.mu.RUnlock()

		data := make(map[string]bool)

		for id := range prov.allocRecvs {
			data[id] = true
		}

		return data
	})

	go simple.StartWatchServer(":8078")
}
