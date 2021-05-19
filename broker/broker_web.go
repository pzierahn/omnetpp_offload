package broker

import (
	"net/http"
)

func (server *broker) pStatusHandle(writer http.ResponseWriter, request *http.Request) {
	server.providers.RLock()
	defer server.providers.RUnlock()

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	//providers := make([]*pb.ProviderState, len(server.providers.provider))
	//
	//inx := 0
	//for _, pro := range server.providers.provider {
	//	providers[inx] = pro
	//	inx++
	//}
	//
	//sort.Slice(providers, func(i, j int) bool {
	//	return providers[i].Utilization.ProviderId < providers[j].Utilization.ProviderId
	//})
	//
	//overview := &pb.ProviderOverview{Providers: providers}
	//byt := simple.PrettyBytes(overview)
	//
	//_, err := writer.Write(byt)
	//if err != nil {
	//	panic(err)
	//}
}

func (server *broker) startWebService() {
	logger.Println("start web service")

	http.HandleFunc("/pstatus", server.pStatusHandle)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
