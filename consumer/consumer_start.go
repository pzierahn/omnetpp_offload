package consumer

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/utils"
	"google.golang.org/grpc"
	"log"
	"path/filepath"
)

func Start(gConf gconfig.GRPCConnection, config *Config) {

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
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	cons := &consumer{
		config: config,
		simulation: &pb.Simulation{
			Id:        simple.NamedId(config.Tag, 8),
			OppConfig: config.OppConfig,
		},
		connections: make(map[string]*connection),
	}
	err = cons.zipSource()
	if err != nil {
		log.Fatalln(err)
	}

	broker := pb.NewBrokerClient(conn)
	providers, err := broker.GetProviders(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("providers: %v", simple.PrettyString(providers.Items))

	for _, prov := range providers.Items {
		_, err = cons.connect(prov)
		if err != nil {
			log.Println(err)
		}
	}

	err = cons.checkoutSimulations()
	if err != nil {
		log.Fatalln(err)
	}

	err = cons.compile()
	if err != nil {
		log.Fatalln(err)
	}

	return
}
