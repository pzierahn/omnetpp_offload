package provider

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/omnetpp"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/sysinfo"
	"path/filepath"
)

func newOpp(simulation *pb.Simulation) (base string, opp omnetpp.OmnetProject) {
	base = filepath.Join(cachePath, simulation.Id)

	conf := &omnetpp.Config{
		OppConfig: simulation.OppConfig,
		Path:      base,
	}

	opp = omnetpp.New(conf)

	return
}

func (prov *provider) compile(ctx context.Context, simulation *pb.Simulation) (bin *pb.Binary, err error) {

	base, opp := newOpp(simulation)
	err = opp.Clean(ctx)
	if err != nil {
		return
	}

	var ref *pb.StorageRef
	var buf bytes.Buffer
	var files map[string]bool
	var cleanFiles map[string]bool
	var buildFiles map[string]bool

	cleanFiles, err = simple.ListDir(base)
	if err != nil {
		return
	}

	err = opp.Setup(ctx, false)
	if err != nil {
		return
	}

	buildFiles, err = simple.ListDir(base)
	if err != nil {
		return
	}

	files = simple.DirDiff(cleanFiles, buildFiles)
	buf, err = simple.TarGzFiles(base, simulation.Id, files)
	if err != nil {
		return
	}

	ref = &pb.StorageRef{
		Bucket:   simulation.Id,
		Filename: fmt.Sprintf("binary/%s.tgz", sysinfo.ArchSignature()),
	}

	err = prov.store.Put(&buf, ref)
	if err != nil {
		return
	}

	bin = &pb.Binary{
		SimulationId: simulation.Id,
		Arch:         sysinfo.Arch(),
		Ref:          ref,
	}

	return
}
