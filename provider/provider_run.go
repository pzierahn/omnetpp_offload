package provider

import (
	"context"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/eval"
	"github.com/pzierahn/omnetpp_offload/omnetpp"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"os"
	"path/filepath"
)

func (prov *provider) run(ctx context.Context, run *pb.SimulationRun) (ref *pb.StorageRef, err error) {

	//
	// Fake copy simulation
	//

	// Simulation directory with simulation source code
	simulationBase := filepath.Join(cachePath, run.SimulationId)

	// This will be the working directory, that contains the results for the job.
	// A fake copy is created to use all configs, ned files and ini files.
	simulationPath := filepath.Join(
		cachePath,
		"fake-copies",
		simple.NamedId(run.SimulationId, 8),
	)

	defer func() {
		// Delete fake-copy after completion
		_ = os.RemoveAll(simulationPath)
	}()

	if err = simple.FakeCopy(simulationBase, simulationPath); err != nil {
		return
	}

	//
	// Execute simulation run
	//

	files := simple.FilesChangeDetector{Root: simulationPath}
	if err = files.Snapshot(); err != nil {
		return
	}

	prov.mu.RLock()
	sess, ok := prov.sessions[run.SimulationId]
	prov.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no session for simulation %s", run.SimulationId)
	}

	oppConf := omnetpp.Config{
		OppConfig: sess.OppConfig,
		Path:      simulationPath,
	}

	done := eval.LogRun(prov.providerId, run.Config, run.RunNum)

	opp := omnetpp.New(&oppConf)
	err = opp.Run(ctx, run.Config, run.RunNum)
	if err != nil {
		return nil, done(err)
	}

	_ = done(nil)

	//
	// Collect and upload results
	//

	filename := fmt.Sprintf("results/%s_%s.tgz", run.Config, run.RunNum)
	done = eval.LogAction(eval.ActionCompress, filename)
	buf, err := files.ZipChanges("")
	_ = done(err)
	if err != nil {
		return
	}

	ref = &pb.StorageRef{
		Bucket:   run.SimulationId,
		Filename: filename,
	}

	err = prov.store.PushFile(&buf, ref)

	return
}
