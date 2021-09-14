package provider

import (
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/eval"
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

	cfiles := simple.ChangedFiles{Root: base}
	if err = cfiles.Snapshot(); err != nil {
		return
	}

	done := eval.LogAction(prov.providerId, eval.ActionCompile)
	err = opp.Setup(ctx, false)
	if err != nil {
		return nil, done(err)
	}

	_ = done(nil)

	buf, err := cfiles.ZipChanges(simulation.Id)
	if err != nil {
		return
	}

	ref := &pb.StorageRef{
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
