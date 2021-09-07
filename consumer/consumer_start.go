package consumer

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/eval"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"google.golang.org/grpc"
	"log"
	"os"
	"path/filepath"
)

func Start(ctx context.Context, config *Config) {

	//go statisticJsonApi()

	if config.Tag == "" {
		config.Tag = filepath.Base(config.Path)
	}

	id := simple.NamedId(config.Tag, 8)
	log.Println("#################################################")
	log.Printf("Start: simulation %s", id)
	log.Println("#################################################")

	eval.ScenarioId = os.Getenv("SCENARIOID")
	eval.TrailId = os.Getenv("TRAILID")
	eval.SimulationId = id

	log.Printf("Start: connecting to broker (%v)", gconfig.BrokerDialAddr())

	conn, err := grpc.Dial(
		gconfig.BrokerDialAddr(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	//log.Printf("Start: set execution timeout to %v", timeout)

	go func() {
		// TODO: Find a more elegant way of doing this
		<-ctx.Done()
		log.Fatalf("Start: execution timeout")
	}()

	cons := &consumer{
		ctx:    ctx,
		config: config,
		bconn:  conn,
		simulation: &pb.Simulation{
			Id:        id,
			OppConfig: config.OppConfig,
		},
		allocate: newQueue(),
	}

	log.Printf("Start: zipping %s", cons.config.Path)

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

	log.Printf("Start: simulation finished!")

	eval.WriteRuns(filepath.Join(config.Path, "opp-edge-eval-runs-"+eval.ScenarioId+".csv"))
	eval.WriteTransfers(filepath.Join(config.Path, "opp-edge-eval-transfers-"+eval.ScenarioId+".csv"))
	eval.WriteSetup(filepath.Join(config.Path, "opp-edge-eval-setup-"+eval.ScenarioId+".csv"))

	return
}
