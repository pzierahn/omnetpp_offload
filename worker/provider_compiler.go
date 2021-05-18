package worker

import (
	"context"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/omnetpp"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"runtime"
)

type Compiler struct {
	Broker         pb.BrokerClient
	Storage        storage.Client
	SimulationId   string
	SimulationBase string
	OppConfig      *pb.OppConfig
	cleanedFiles   map[string]bool
	compiledFiles  map[string]bool
}

func (compiler Compiler) Opp() (opp omnetpp.OmnetProject) {
	opp = omnetpp.New(&omnetpp.Config{
		OppConfig: compiler.OppConfig,
		Path:      compiler.SimulationBase,
	})

	return
}

func (compiler Compiler) CheckinBinary() (err error) {

	buf, err := simple.TarGzFiles(
		compiler.SimulationBase,
		compiler.SimulationId,
		compiler.compiledFiles)
	if err != nil {
		return
	}

	arch := &pb.OsArch{
		Os:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
	ref, err := compiler.Storage.Upload(&buf, storage.FileMeta{
		Bucket:   compiler.SimulationId,
		Filename: fmt.Sprintf("binary_%s_%s.tar.gz", arch.Os, arch.Arch),
	})
	if err != nil {
		return
	}

	_, err = compiler.Broker.AddBinary(context.Background(), &pb.SimBinary{
		SimulationId: compiler.SimulationId,
		Arch:         arch,
		Binary:       ref,
	})
	if err != nil {
		_, _ = compiler.Storage.Delete(ref)
		return
	}

	return
}

func (compiler Compiler) Clean() (err error) {

	opp := compiler.Opp()
	err = opp.Clean()
	if err != nil {
		return
	}

	compiler.cleanedFiles, err = simple.ListDir(opp.Path)

	return
}

func (compiler Compiler) Compile() (err error) {

	opp := compiler.Opp()
	err = opp.Setup(false)
	if err != nil {
		return
	}

	compiler.compiledFiles, err = simple.ListDir(opp.Path)
	if err != nil {
		return
	}

	return
}

func (compiler Compiler) CompiledFiles() (files map[string]bool) {
	files = simple.DirDiff(compiler.cleanedFiles, compiler.compiledFiles)
	return
}
