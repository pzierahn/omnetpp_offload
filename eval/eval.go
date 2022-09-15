package eval

import pb "github.com/pzierahn/omnetpp_offload/proto"

const (
	_ = uint32(iota)
	StateStarted
	StateFinished
	StateFailed
)

const (
	ActivityCompile  = "Compile"
	ActivityRun      = "Run"
	ActivityUpload   = "Upload"
	ActivityDownload = "Download"
	ActivityCompress = "Compress"
	ActivityExtract  = "Extract"
)

type Event struct {
	Activity      string
	SimulationRun *pb.SimulationRun
	Filename      string
}

func (event Event) runId() (conf string, num string) {
	if event.SimulationRun == nil {
		return "", ""
	}

	return event.SimulationRun.Config, event.SimulationRun.RunNum
}
