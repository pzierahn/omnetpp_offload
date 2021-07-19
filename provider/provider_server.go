package provider

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"github.com/pzierahn/project.go.omnetpp/sysinfo"
	"google.golang.org/grpc"
	"log"
	"runtime"
	"sync"
	"time"
)

func Start() {

	store := &storage.Server{}

	prov := &provider{
		providerId: simple.NamedId(gconfig.Config.Worker.Name, 8),
		store:      store,
		// TODO: replace this with gconfig values!
		slots:       uint32(runtime.NumCPU()),
		freeSlots:   int32(runtime.NumCPU()),
		cond:        sync.NewCond(&sync.Mutex{}),
		requests:    make(map[simulationId]uint32),
		assignments: make(map[simulationId]uint32),
		allocate:    make(map[simulationId]chan<- uint32),
	}

	log.Printf("start provider (%v)", prov.providerId)

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
	// Start provider
	//

	go prov.listenP2P()
	go prov.listenRelay(brokerConn)
	prov.allocator()

	return
}
