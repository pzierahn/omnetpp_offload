package consumer

import (
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"google.golang.org/grpc"
	"log"
	"path/filepath"
)

func Start(config *Config) {

	go statisticJsonApi()

	if config.Tag == "" {
		config.Tag = filepath.Base(config.Path)
	}

	log.Printf("connecting to broker (%v)", gconfig.BrokerDialAddr())

	conn, err := grpc.Dial(
		gconfig.BrokerDialAddr(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	cons := &consumer{
		config: config,
		bconn:  conn,
		simulation: &pb.Simulation{
			Id:        simple.NamedId(config.Tag, 8),
			OppConfig: config.OppConfig,
		},
		allocate: newQueue(),
	}

	log.Printf("zipping simulation source: %s", cons.config.Path)

	buf, err := simple.TarGz(cons.config.Path, cons.simulation.Id, cons.config.Ignore...)
	if err != nil {
		log.Fatalln(err)
	}

	cons.simulationSource = buf.Bytes()

	onInit := make(chan int32)
	defer close(onInit)

	go cons.startConnector(onInit)

	cons.finished.Add(int(<-onInit))
	cons.finished.Wait()

	log.Printf("simulation finished!")

	showStatistic()

	return
}
