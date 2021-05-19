package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/utils"
	"net/http"
	"sort"
)

func (server *broker) pStatusHandle(writer http.ResponseWriter, request *http.Request) {
	server.providers.RLock()
	defer server.providers.RUnlock()

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	query := request.URL.Query()

	proto := utils.QueryBool(query, "proto", false)

	providers := make([]*pb.ProviderState, len(server.providers.provider))

	inx := 0
	for _, pro := range server.providers.provider {
		providers[inx] = &pb.ProviderState{
			ProviderId:  pro.id,
			Arch:        pro.arch,
			NumCPUs:     pro.numCPUs,
			CpuUsage:    pro.utilization.CpuUsage,
			MemoryUsage: pro.utilization.MemoryUsage,
			Updated:     pro.utilization.Updated,
			Assignments: pro.assignments,
		}
		inx++
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].ProviderId < providers[j].ProviderId
	})

	overview := &pb.ProviderOverview{Items: providers}
	utils.Response(writer, overview, proto)
}

func (server *broker) startWebService() {
	logger.Println("start web service")

	http.HandleFunc("/pstatus", server.pStatusHandle)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
