package compile

import (
	"context"
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"runtime"
)

func (compiler Compiler) Checkin() (err error) {

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
