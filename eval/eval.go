package eval

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
