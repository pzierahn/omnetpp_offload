package provider

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/egrpc"
	"github.com/pzierahn/project.go.omnetpp/eval"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"github.com/pzierahn/project.go.omnetpp/sysinfo"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type provider struct {
	pb.UnimplementedProviderServer
	providerId  string
	store       *storage.Server
	mu          *sync.RWMutex
	cond        *sync.Cond
	slots       uint32
	freeSlots   int32
	requests    map[simulationId]uint32
	assignments map[simulationId]uint32
	allocate    map[simulationId]chan<- uint32
	sessions    map[simulationId]*pb.Session
}

func Start() {

	store := &storage.Server{}

	mu := &sync.RWMutex{}
	prov := &provider{
		providerId:  simple.NamedId(gconfig.Config.Worker.Name, 8),
		store:       store,
		slots:       uint32(gconfig.Jobs()),
		freeSlots:   int32(gconfig.Jobs()),
		mu:          mu,
		cond:        sync.NewCond(mu),
		requests:    make(map[simulationId]uint32),
		assignments: make(map[simulationId]uint32),
		allocate:    make(map[simulationId]chan<- uint32),
		sessions:    make(map[simulationId]*pb.Session),
	}

	log.Printf("start provider (%v)", prov.providerId)

	//
	// Init stuff
	//

	eval.DeviceId = prov.providerId
	prov.recoverSessions()

	//
	// Register provider
	//

	log.Printf("connect to broker %v", gconfig.BrokerDialAddr())

	brokerConn, err := grpc.Dial(
		gconfig.BrokerDialAddr(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	eval.Init(brokerConn)

	broker := pb.NewBrokerClient(brokerConn)

	stream, err := broker.Register(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	err = stream.Send(&pb.Ping{Cast: &pb.Ping_Register{Register: prov.info()}})
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for range time.Tick(time.Millisecond * 500) {

			var util *pb.Utilization
			util, err = sysinfo.GetUtilization(context.Background())
			if err != nil {
				log.Fatalln(err)
			}

			//log.Printf("Start: send utilization %v", util.CpuUsage)

			err = stream.Send(&pb.Ping{Cast: &pb.Ping_Util{Util: util}})
			if err != nil {
				// TODO: reconnect after EOF
				log.Fatalln(err)
			}
		}
	}()

	//
	// Start gRPC server.
	//

	server := grpc.NewServer()
	pb.RegisterProviderServer(server, prov)
	pb.RegisterStorageServer(server, prov.store)

	go egrpc.ServeLocal(prov.providerId, server)
	go egrpc.ServeP2P(prov.providerId, server)
	go egrpc.ServeRelay(prov.providerId, server)

	//
	// Start resource allocator.
	//

	prov.startAllocator()

	return
}
