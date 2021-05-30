package broker

import "log"

//import (
//	pb "github.com/pzierahn/project.go.omnetpp/proto"
//	"github.com/pzierahn/project.go.omnetpp/utils"
//	"google.golang.org/protobuf/proto"
//	"net/http"
//	"sort"
//)
//
//func (server *broker) pStatusHandle(writer http.ResponseWriter, request *http.Request) {
//	server.providers.RLock()
//	defer server.providers.RUnlock()
//
//	writer.Header().Set("Content-Type", "application/json")
//	writer.Header().Set("Access-Control-Allow-Origin", "*")
//
//	query := request.URL.Query()
//	sendProto := utils.QueryBool(query, "proto", false)
//
//	providers := make([]*pb.ProviderState, len(server.providers.provider))
//
//	inx := 0
//
//	for _, prov := range server.providers.provider {
//		prov.RLock()
//
//		assignments := make(map[string]*pb.SimulationRun)
//
//		for key, val := range prov.assignments {
//			assignments[string(key)] = proto.Clone(val).(*pb.SimulationRun)
//		}
//
//		providers[inx] = &pb.ProviderState{
//			ProviderId:  prov.id,
//			Arch:        prov.arch,
//			NumCPUs:     prov.numCPUs,
//			Utilization: proto.Clone(prov.utilization).(*pb.Utilization),
//			Assignments: assignments,
//			Building:    prov.building,
//		}
//		prov.RUnlock()
//
//		inx++
//	}
//
//	sort.Slice(providers, func(i, j int) bool {
//		return providers[i].ProviderId < providers[j].ProviderId
//	})
//
//	overview := &pb.ProviderOverview{Items: providers}
//	utils.Response(writer, overview, sendProto)
//}
//
//func (server *broker) sStatusHandle(writer http.ResponseWriter, request *http.Request) {
//	server.providers.RLock()
//	defer server.providers.RUnlock()
//
//	writer.Header().Set("Content-Type", "application/json")
//	writer.Header().Set("Access-Control-Allow-Origin", "*")
//
//	query := request.URL.Query()
//	sendProto := utils.QueryBool(query, "proto", false)
//
//	var sim *simulationState
//
//	server.simulations.RLock()
//	for _, sim = range server.simulations.simulations {
//		break
//	}
//	server.simulations.RUnlock()
//
//	overview := &pb.SimulationState{}
//
//	if sim != nil {
//		sim.RLock()
//
//		runs := make(map[string]*pb.SimulationRun)
//		for id, run := range sim.runs {
//			runs[string(id)] = proto.Clone(run).(*pb.SimulationRun)
//		}
//
//		queue := make(map[string]bool)
//		for id, val := range sim.queue {
//			queue[string(id)] = val
//		}
//
//		binaries := make(map[string]*pb.Binary)
//		for id, bin := range sim.binaries {
//			binaries[string(id)] = proto.Clone(bin).(*pb.Binary)
//		}
//
//		overview = &pb.SimulationState{
//			SimulationId: sim.simulationId,
//			Queue:        queue,
//			Runs:         runs,
//			Source:       proto.Clone(sim.source).(*pb.StorageRef),
//			OppConfig:    proto.Clone(sim.oppConfig).(*pb.OppConfig),
//			Binaries:     binaries,
//		}
//
//		sim.RUnlock()
//	}
//
//	utils.Response(writer, overview, sendProto)
//}

func (broker *broker) startWebService() {

	log.Printf("startWebService STUB!")

	//logger.Println("start web service on http://localhost:8090")
	//
	//http.HandleFunc("/provider", server.pStatusHandle)
	//http.HandleFunc("/simulation", server.sStatusHandle)
	//
	//err := http.ListenAndServe(":8090", nil)
	//if err != nil {
	//	panic(err)
	//}
}
