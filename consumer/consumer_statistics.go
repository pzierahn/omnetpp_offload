package consumer

import (
	"encoding/json"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

type ProviderStatistic struct {
	ProviderId       string
	ProviderInfo     *pb.ProviderInfo
	ExecutedTasks    int
	ExecutionTimeSum string
	ExecutionTimeAvg string
}

var execMu sync.RWMutex
var providerInfos = make(map[string]*pb.ProviderInfo)

type durationStatistic map[string][]time.Duration

var execTime = make(durationStatistic)

func logProviderInfo(providerId string, info *pb.ProviderInfo) {
	execMu.Lock()
	defer execMu.Unlock()

	providerInfos[providerId] = info
}

func logExecTime(providerId string, dur time.Duration) {
	execMu.Lock()
	defer execMu.Unlock()

	execTime[providerId] = append(execTime[providerId], dur)
}

func handleStatistic(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	if bytes, err := json.MarshalIndent(gatherStatistic(), "", "  "); err != nil {
		writer.WriteHeader(503)
		_, _ = fmt.Fprint(writer, err.Error())
		log.Println(err)
	} else {
		_, _ = writer.Write(bytes)
	}
}

func statisticJsonApi() {
	log.Println("start web service on http://localhost:8800/")

	http.HandleFunc("/", handleStatistic)

	err := http.ListenAndServe(":8800", nil)
	if err != nil {
		panic(err)
	}
}

func gatherStatistic() (stats []ProviderStatistic) {

	execMu.RLock()
	defer execMu.RUnlock()

	stats = make([]ProviderStatistic, 0)

	for prov, exeTimes := range execTime {
		var set time.Duration

		for _, dur := range exeTimes {
			set += dur
		}

		stat := ProviderStatistic{
			ProviderId:       prov,
			ProviderInfo:     providerInfos[prov],
			ExecutedTasks:    len(exeTimes),
			ExecutionTimeSum: fmt.Sprintf("%v", set),
			ExecutionTimeAvg: fmt.Sprintf("%v", time.Duration(set.Nanoseconds()/int64(len(exeTimes)))),
		}

		stats = append(stats, stat)
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].ProviderId < stats[j].ProviderId
	})

	return
}

func showStatistic() {
	log.Println("execution statistics", simple.PrettyString(gatherStatistic()))
}
