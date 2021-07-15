package provider

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"github.com/pzierahn/project.go.omnetpp/sysinfo"
	"google.golang.org/grpc"
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

func Start(conf gconfig.Config) {

	store := &storage.Server{}

	prov := &provider{
		providerId: simple.NamedId(conf.Worker.Name, 8),
		store:      store,
		slots:      uint32(runtime.NumCPU()),
		freeSlots:  int32(runtime.NumCPU()),
		//slots:       uint32(1),
		//freeSlots:   1,
		cond:        sync.NewCond(&sync.Mutex{}),
		requests:    make(map[simulationId]uint32),
		assignments: make(map[simulationId]uint32),
		runCtx:      make(map[simulationId]context.CancelFunc),
		allocate:    make(map[simulationId]chan<- uint32),
	}

	log.Printf("start provider (%v)", prov.providerId)

	//
	// Register provider
	//

	brokerConn, err := grpc.Dial(
		conf.Broker.DialAddr(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalln(err)
	}

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
			util, err = sysinfo.GetUtilization()
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
	// Start provider
	//

	go prov.allocator()

	for {
		log.Println("wait for peer to peer connect")

		p2p, err := stargate.DialQUICgRPCListener(context.Background(), prov.providerId)
		if err != nil {
			log.Println(err)
			continue
		}

		go func(p2p net.Listener) {

			defer func() { _ = p2p.Close() }()

			server := grpc.NewServer()
			pb.RegisterProviderServer(server, prov)
			pb.RegisterStorageServer(server, store)
			err := server.Serve(p2p)
			if err != nil {
				log.Println(err)
			}
		}(p2p)
	}

	return
}
