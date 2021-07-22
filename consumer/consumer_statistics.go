package consumer

import (
	"encoding/json"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/statistic"
	"log"
	"net/http"
	"sort"
	"sync"
)

type Statistic struct {
	Info      *pb.ProviderInfo
	Execution *statistic.TimeStatistic
	Checkout  *statistic.TimeStatistic
	Upload    *statistic.LoadStatistic
	Download  *statistic.LoadStatistic
	Compile   *statistic.TimeStatistic
}

type pstatistic struct {
	mu        sync.RWMutex
	Info      map[string]*pb.ProviderInfo
	Execution map[string]*statistic.Time
	Checkout  map[string]*statistic.Time
	Upload    map[string]*statistic.Load
	Download  map[string]*statistic.Load
	Compile   map[string]*statistic.Time
}

func (stat *pstatistic) SetInfo(id string, info *pb.ProviderInfo) {
	stat.mu.Lock()
	defer stat.mu.Unlock()

	stat.Info[id] = info
}

func (stat *pstatistic) GetExecution(id string) (exec *statistic.Time) {
	stat.mu.Lock()
	defer stat.mu.Unlock()

	var ok bool
	if exec, ok = stat.Execution[id]; ok {
		return
	}

	exec = &statistic.Time{}
	stat.Execution[id] = exec

	return
}

func (stat *pstatistic) GetDownload(id string) (load *statistic.Load) {
	stat.mu.Lock()
	defer stat.mu.Unlock()

	var ok bool
	if load, ok = stat.Download[id]; ok {
		return
	}

	load = &statistic.Load{}
	stat.Download[id] = load

	return
}

func (stat *pstatistic) GetUpload(id string) (load *statistic.Load) {
	stat.mu.Lock()
	defer stat.mu.Unlock()

	var ok bool
	if load, ok = stat.Upload[id]; ok {
		return
	}

	load = &statistic.Load{}
	stat.Upload[id] = load

	return
}

func (stat *pstatistic) GetCheckout(id string) (check *statistic.Time) {
	stat.mu.Lock()
	defer stat.mu.Unlock()

	var ok bool
	if check, ok = stat.Checkout[id]; ok {
		return
	}

	check = &statistic.Time{}
	stat.Checkout[id] = check

	return
}

func (stat *pstatistic) GetCompile(id string) (check *statistic.Time) {
	stat.mu.Lock()
	defer stat.mu.Unlock()

	var ok bool
	if check, ok = stat.Compile[id]; ok {
		return
	}

	check = &statistic.Time{}
	stat.Compile[id] = check

	return
}

func (stat *pstatistic) Export() (list []Statistic) {
	stat.mu.Lock()
	defer stat.mu.Unlock()

	list = make([]Statistic, 0)

	for id := range stat.Info {
		item := Statistic{
			Info:      stat.Info[id],
			Execution: stat.Execution[id].Export(),
			Checkout:  stat.Checkout[id].Export(),
			Upload:    stat.Upload[id].Export(),
			Download:  stat.Download[id].Export(),
			Compile:   stat.Compile[id].Export(),
		}

		list = append(list, item)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Info.ProviderId < list[j].Info.ProviderId
	})

	return
}

var stat = &pstatistic{
	Info:      make(map[string]*pb.ProviderInfo),
	Execution: make(map[string]*statistic.Time),
	Checkout:  make(map[string]*statistic.Time),
	Upload:    make(map[string]*statistic.Load),
	Download:  make(map[string]*statistic.Load),
	Compile:   make(map[string]*statistic.Time),
}

func handleStatistic(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	if bytes, err := json.MarshalIndent(stat.Export(), "", "  "); err != nil {
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

func showStatistic() {
	log.Println("execution statistics", simple.PrettyString(stat.Export()))
}
