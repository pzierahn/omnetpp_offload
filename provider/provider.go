package provider

import (
	"context"
	"github.com/pzierahn/omnetpp_offload/eval"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"github.com/pzierahn/omnetpp_offload/stargate"
	"github.com/pzierahn/omnetpp_offload/stargrpc"
	"github.com/pzierahn/omnetpp_offload/storage"
	"github.com/pzierahn/omnetpp_offload/sysinfo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

type simulationId = string

type provider struct {
	pb.UnimplementedProviderServer
	providerId     string
	numJobs        int
	store          *storage.Server
	slots          chan int
	mu             *sync.RWMutex
	sessions       map[simulationId]*pb.Session
	executionTimes map[simulationId]time.Duration
	newRecv        *sync.Cond
	allocRecvs     map[simulationId]chan<- int
}

func Start(config gconfig.Config) {

	mu := &sync.RWMutex{}
	prov := &provider{
		providerId:     simple.NamedId(config.Provider.Name, 8),
		numJobs:        config.Provider.Jobs,
		store:          &storage.Server{},
		slots:          make(chan int, config.Provider.Jobs),
		mu:             mu,
		newRecv:        sync.NewCond(mu),
		sessions:       make(map[simulationId]*pb.Session),
		executionTimes: make(map[simulationId]time.Duration),
		allocRecvs:     make(map[simulationId]chan<- int),
	}

	log.Printf("start provider (%v)", prov.providerId)

	//
	// Init stuff
	//

	prov.recoverSessions()

	startWatchers(prov)

	//
	// Register provider
	//

	log.Printf("connect to broker %v", config.Broker.BrokerDialAddr())

	brokerConn, err := grpc.Dial(
		config.Broker.BrokerDialAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		_ = brokerConn.Close()
	}()

	eval.Init(brokerConn, prov.providerId)

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
	// Start stargate-gRPC servers.
	//

	server := grpc.NewServer()
	pb.RegisterProviderServer(server, prov)
	pb.RegisterStorageServer(server, prov.store)

	stargate.SetConfig(stargate.Config{
		Addr: config.Broker.Address,
		Port: config.Broker.StargatePort,
	})

	go stargrpc.ServeLocal(prov.providerId, server)
	go stargrpc.ServeP2P(prov.providerId, server)
	go stargrpc.ServeRelay(prov.providerId, server)

	//
	// Start resource allocator.
	//

	prov.startAllocator(config.Provider.Jobs)

	return
}
