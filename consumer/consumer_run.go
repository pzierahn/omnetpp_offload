package consumer

import (
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"github.com/pzierahn/project.go.omnetpp/utils"
	"google.golang.org/grpc"
	"log"
	"path/filepath"
	"time"
)

func Run(gConf gconfig.GRPCConnection, config *Config) (err error) {

	if config.Tag == "" {
		config.Tag = filepath.Base(config.Path)
	}

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

		pconn, remote, err := stargate.Dial(ctx, prov.ProviderId)
		if err != nil {
			// Connection failed!
			log.Println(err)
			continue
		}

		log.Printf("connected to %v", remote)

		qconn, err := grpc.Dial(
			remote.String(),
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithContextDialer(utils.GRPCDialer(pconn)),
		)
		if err != nil {
			log.Fatalln(err)
		}

		provider := pb.NewProviderClient(qconn)

		start := time.Now()

		for inx := 0; inx < 5000; inx++ {
			if inx%100 == 0 {
				log.Printf("request: %v", inx)
			}

			_, err = provider.Info(context.Background(), &pb.Empty{})
			if err != nil {
				log.Fatalln(err)
			}
		}

		end := time.Now()

		log.Printf("exectime: %v", end.Sub(start))
		log.Printf("average exectime: %v", end.Sub(start)/10_000)

		_ = pconn.Close()
	}

	return
}
