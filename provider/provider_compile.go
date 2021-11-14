package provider

import (
	"context"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/eval"
	"github.com/pzierahn/omnetpp_offload/omnetpp"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"github.com/pzierahn/omnetpp_offload/sysinfo"
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
	err = opp.MakeMake(ctx)
	if err != nil {
		return
	}

	err = opp.Clean(ctx)
	if err != nil {
		return
	}

	cfiles := simple.FilesChangeDetector{Root: base}
	if err = cfiles.Snapshot(); err != nil {
		return
	}

	done := eval.LogAction(eval.ActionCompile, sysinfo.ArchSignature())
	err = opp.Setup(ctx, false)
	_ = done(nil)
	if err != nil {
		return
	}

	filename := fmt.Sprintf("binary/%s.tgz", sysinfo.ArchSignature())
	done = eval.LogAction(eval.ActionCompress, filename)
	buf, err := cfiles.ZipChanges(simulation.Id)
	_ = done(err)
	if err != nil {
		return
	}

	ref := &pb.StorageRef{
		Bucket:   simulation.Id,
		Filename: filename,
	}

	err = prov.store.PushFile(&buf, ref)
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
