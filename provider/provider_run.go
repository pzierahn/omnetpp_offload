package provider

import (
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/eval"
	"github.com/pzierahn/project.go.omnetpp/omnetpp"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"os"
	"path/filepath"
)

func (prov *provider) run(ctx context.Context, run *pb.SimulationRun) (ref *pb.StorageRef, err error) {

	//
	// Setup mirror simulation
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

	files := simple.ChangedFiles{Root: simulationPath}
	if err = files.Snapshot(); err != nil {
		return
	}

	oppConf := omnetpp.Config{
		OppConfig: run.OppConfig,
		Path:      simulationPath,
	}

	done := eval.LogRun(prov.providerId, run.Config, run.RunNum)

	opp := omnetpp.New(&oppConf)
	err = opp.RunContext(ctx, run.Config, run.RunNum)
	if err != nil {
		return nil, done(err)
	}

	_ = done(nil)

	//
	// Collect and upload results
	//

	buf, err := files.ZipChanges("")
	if err != nil {
		return
	}

	ref = &pb.StorageRef{
		Bucket:   run.SimulationId,
		Filename: fmt.Sprintf("results/%s_%s.tgz", run.Config, run.RunNum),
	}

	err = prov.store.Put(&buf, ref)

	return
}
