package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/utils"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sort"
)

func (server *broker) pStatusHandle(writer http.ResponseWriter, request *http.Request) {
	server.providers.RLock()
	defer server.providers.RUnlock()

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	query := request.URL.Query()
	sendProto := utils.QueryBool(query, "proto", false)

	providers := make([]*pb.ProviderState, len(server.providers.provider))

	inx := 0

	server.providers.RLock()
	for _, pro := range server.providers.provider {
		pro.RLock()

		assignments := make(map[string]*pb.Assignment)

		for key, val := range pro.assignments {
			assignments[key] = proto.Clone(val).(*pb.Assignment)
		}

		providers[inx] = &pb.ProviderState{
			ProviderId:  pro.id,
			Arch:        pro.arch,
			NumCPUs:     pro.numCPUs,
			Utilization: proto.Clone(pro.utilization).(*pb.Utilization),
			Assignments: assignments,
		}
		pro.RUnlock()

		inx++
	}
	server.providers.RUnlock()

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].ProviderId < providers[j].ProviderId
	})

	overview := &pb.ProviderOverview{Items: providers}
	utils.Response(writer, overview, sendProto)
}

func (server *broker) sStatusHandle(writer http.ResponseWriter, request *http.Request) {
	server.providers.RLock()
	defer server.providers.RUnlock()

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	query := request.URL.Query()
	sendProto := utils.QueryBool(query, "proto", false)

	var sim *simulationState

	server.simulations.RLock()
	for _, sim = range server.simulations.simulations {
		break
	}
	server.simulations.RUnlock()

	overview := &pb.SimulationState{}

	if sim != nil {
		sim.RLock()

		runs := make(map[string]*pb.SimulationRun)
		for id, run := range sim.runs {
			runs[string(id)] = proto.Clone(run).(*pb.SimulationRun)
		}

		queue := make(map[string]bool)
		for id, val := range sim.queue {
			queue[string(id)] = val
		}

		binaries := make(map[string]*pb.Binary)
		for id, bin := range sim.binaries {
			binaries[string(id)] = proto.Clone(bin).(*pb.Binary)
		}

		overview = &pb.SimulationState{
			SimulationId: sim.simulationId,
			Queue:        queue,
			Runs:         runs,
			Source:       proto.Clone(sim.source).(*pb.StorageRef),
			OppConfig:    proto.Clone(sim.oppConfig).(*pb.OppConfig),
			Binaries:     binaries,
		}

		sim.RUnlock()
	}

	utils.Response(writer, overview, sendProto)
}

func (server *broker) startWebService() {
	logger.Println("start web service")

	http.HandleFunc("/provider", server.pStatusHandle)
	http.HandleFunc("/simulation", server.sStatusHandle)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
