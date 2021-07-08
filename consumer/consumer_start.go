package consumer

import (
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/equic"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"path/filepath"
	"sync"
)

func Start(gConf gconfig.GRPCConnection, config *Config) {

	if config.Tag == "" {
		config.Tag = filepath.Base(config.Path)
	}

	log.Printf("connecting to broker (%v)", gConf.DialAddr())

	_, dialer := equic.GRPCDialerAuto()
	conn, err := grpc.Dial(
		gConf.DialAddr(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(dialer),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	cons := &consumer{
		consumerId: fmt.Sprintf("consumer-%x", rand.Uint32()),
		config:     config,
		simulation: &pb.Simulation{
			Id:        simple.NamedId(config.Tag, 8),
			OppConfig: config.OppConfig,
		},
		connections: make(map[string]*connection),
		allocCond:   sync.NewCond(&sync.Mutex{}),
		allocator:   make(chan *pb.SimulationRun),
	}

	log.Printf("zipping simulation source: %s", cons.config.Path)

	buf, err := simple.TarGz(cons.config.Path, cons.simulation.Id, cons.config.Ignore...)
	if err != nil {
		log.Fatalln(err)
	}

	cons.simulationSource = buf.Bytes()

	broker := pb.NewBrokerClient(conn)
	go cons.startConnector(broker)

	err = cons.dispatchTasks()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("simulation finished!")

	return
}
