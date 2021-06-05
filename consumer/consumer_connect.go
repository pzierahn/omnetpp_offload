package consumer

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"github.com/pzierahn/project.go.omnetpp/utils"
	"google.golang.org/grpc"
	"log"
	"time"
)

type connection struct {
	info     *pb.ProviderInfo
	provider pb.ProviderClient
	store    pb.StorageClient
}

func (cons *consumer) connect(prov *pb.ProviderInfo) (conn *connection, err error) {

	log.Printf("connect to provider (%v)", prov.ProviderId)

	ctx, cln := context.WithTimeout(context.Background(), time.Second*4)
	defer cln()

	gate, remote, err := stargate.Dial(ctx, prov.ProviderId)
	if err != nil {
		// Connection failed!
		return
	}

	var cConn *grpc.ClientConn
	cConn, err = grpc.Dial(
		remote.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(utils.GRPCDialer(gate)),
	)
	if err != nil {
		return
	}

	conn = &connection{
		info:     prov,
		provider: pb.NewProviderClient(cConn),
		store:    pb.NewStorageClient(cConn),
	}

	cons.connMu.Lock()
	cons.connections[prov.ProviderId] = conn
	cons.connMu.Unlock()

	// TODO: Handle disconnect!

	//storeCli := storage.FromClient(store)
	//
	//log.Println("zipping", config.Path)
	//
	//var buf bytes.Buffer
	//buf, err = simple.TarGz(config.Path, cons.simulationId, config.Ignore...)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Printf("upload: %s (%d bytes)", cons.simulationId, buf.Len())
	//
	//var ref *pb.StorageRef
	//ref, err = storeCli.Upload(&buf, storage.FileMeta{
	//	Bucket:   cons.simulationId,
	//	Filename: "zipSource.tar.gz",
	//})
	//
	//log.Printf("checkout: %v", cons.simulationId)
	//
	//_, err = provider.Checkout(context.Background(), &pb.Bundle{
	//	SimulationId: cons.simulationId,
	//	Source:       ref,
	//})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Printf("compile: %v", cons.simulationId)
	//
	//_, err = provider.Compile(context.Background(), &pb.Simulation{
	//	Id:        cons.simulationId,
	//	OppConfig: config.OppConfig,
	//})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Printf("list run numbers: %v", cons.simulationId)
	//
	//runNums, err := provider.ListRunNums(context.Background(), &pb.Simulation{
	//	Id:        cons.simulationId,
	//	OppConfig: config.OppConfig,
	//	Config:    config.SimulateConfigs[0],
	//})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Printf("runNums: %v", runNums)
	//
	//var wg sync.WaitGroup
	//work := make(chan *pb.Simulation)
	//
	//for inx := 0; inx < 4; inx++ {
	//	go func(agent int) {
	//		for sim := range work {
	//			log.Printf("[%d] run: %v", agent, sim.RunNum)
	//
	//			ref, err = provider.Run(context.Background(), sim)
	//			if err != nil {
	//				log.Fatalln(err)
	//			}
	//
	//			log.Printf("[%d] result: %v", agent, ref)
	//			wg.Done()
	//		}
	//	}(inx)
	//}
	//
	//for _, num := range runNums.Runs {
	//	wg.Add(1)
	//	work <- &pb.Simulation{
	//		Id:        cons.simulationId,
	//		OppConfig: config.OppConfig,
	//		Config:    runNums.Config,
	//		RunNum:    num,
	//	}
	//}
	//
	//wg.Wait()

	//err = pconn.Close()

	return
}
