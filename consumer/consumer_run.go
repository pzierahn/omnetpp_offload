package consumer

import (
	"context"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/stargate"
	"github.com/patrickz98/project.go.omnetpp/utils"
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

	ctx := context.Background()

	broker := pb.NewBrokerClient(conn)
	providers, err := broker.GetProviders(ctx, &pb.Empty{})
	if err != nil {
		return
	}

	log.Printf("providers %d", len(providers.Items))

	for _, prov := range providers.Items {

		log.Printf("connect to provider (%v)", prov.ProviderId)

		pconn, remote := stargate.Connect(prov.ProviderId)

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

		for range time.Tick(time.Second) {
			status, err := provider.Status(context.Background(), &pb.Empty{})
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("provider status (%s) %v", prov.ProviderId, status)
		}

		_ = pconn.Close()
	}

	return
}
