package consumer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"github.com/pzierahn/project.go.omnetpp/storage"
	"github.com/pzierahn/project.go.omnetpp/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"path/filepath"
	"sync"
	"time"
)

func Run(gConf gconfig.GRPCConnection, config *Config) (err error) {

	if config.Tag == "" {
		config.Tag = filepath.Base(config.Path)
	}

	simulationId := simple.NamedId(config.Tag, 8)

	log.Printf("connecting to broker (%v)", gConf.DialAddr())

	_, dialer := utils.GRPCDialerAuto()
	conn, err := grpc.Dial(
		gConf.DialAddr(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithContextDialer(dialer),
	)
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	broker := pb.NewBrokerClient(conn)
	providers, err := broker.GetProviders(context.Background(), &pb.Empty{})
	if err != nil {
		return
	}

	log.Printf("providers %v", simple.PrettyString(providers.Items))

	for _, prov := range providers.Items {

		log.Printf("connect to provider (%v)", prov.ProviderId)

		ctx, _ := context.WithTimeout(context.Background(), time.Second*4)

		var pconn *net.UDPConn
		var remote *net.UDPAddr

		pconn, remote, err = stargate.Dial(ctx, prov.ProviderId)
		if err != nil {
			// Connection failed!
			log.Println(err)
			continue
		}

		log.Printf("connected to %v", remote)

		var qconn *grpc.ClientConn
		qconn, err = grpc.Dial(
			remote.String(),
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithContextDialer(utils.GRPCDialer(pconn)),
		)
		if err != nil {
			log.Fatalln(err)
		}

		provider := pb.NewProviderClient(qconn)
		store := pb.NewStorageClient(qconn)
		storeCli := storage.FromClient(store)

		log.Println("zipping", config.Path)

		var buf bytes.Buffer
		buf, err = simple.TarGz(config.Path, simulationId, config.Exclude...)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("upload: %s (%d bytes)", simulationId, buf.Len())

		var ref *pb.StorageRef
		ref, err = storeCli.Upload(&buf, storage.FileMeta{
			Bucket:   simulationId,
			Filename: "source.tar.gz",
		})

		log.Printf("checkout: %v", simulationId)

		_, err = provider.Checkout(context.Background(), &pb.Bundle{
			SimulationId: simulationId,
			Source:       ref,
		})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("compile: %v", simulationId)

		_, err = provider.Compile(context.Background(), &pb.Simulation{
			Id:        simulationId,
			OppConfig: config.OppConfig,
		})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("list run numbers: %v", simulationId)

		runNums, err := provider.ListRunNums(context.Background(), &pb.Simulation{
			Id:        simulationId,
			OppConfig: config.OppConfig,
			Config:    config.SimulateConfigs[0],
		})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("runNums: %v", runNums)

		var wg sync.WaitGroup
		work := make(chan *pb.Simulation)

		for inx := 0; inx < 4; inx++ {
			go func(agent int) {
				for sim := range work {
					log.Printf("[%d] run: %v", agent, sim.RunNum)

					ref, err = provider.Run(context.Background(), sim)
					if err != nil {
						log.Fatalln(err)
					}

					log.Printf("[%d] result: %v", agent, ref)
					wg.Done()
				}
			}(inx)
		}

		for _, num := range runNums.Runs {
			wg.Add(1)
			work <- &pb.Simulation{
				Id:        simulationId,
				OppConfig: config.OppConfig,
				Config:    runNums.Config,
				RunNum:    num,
			}
		}

		wg.Wait()

		_ = pconn.Close()
	}

	return
}
