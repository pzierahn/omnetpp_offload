package broker

import (
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/utils"
	"log"
	"net/http"
)

func (broker *broker) pStatusHandle(writer http.ResponseWriter, request *http.Request) {
	broker.pmu.RLock()
	defer broker.pmu.RUnlock()

	writer.Header().Set("Access-Control-Allow-Origin", "*")

	query := request.URL.Query()
	sendProto := utils.QueryBool(query, "proto", false)

	loads := make(map[string]*pb.Utilization)

	broker.umu.RLock()
	for id, utilization := range broker.utilization {
		loads[id] = utilization
	}
	broker.umu.RUnlock()

	overview := &pb.Utilizations{Providers: loads}
	utils.Response(writer, overview, sendProto)
}

func (broker *broker) startWebService() {

	log.Println("start web service on http://localhost:8090/providers")

	http.HandleFunc("/providers", broker.pStatusHandle)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
