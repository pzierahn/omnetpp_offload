package simple

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

func CheckError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// ReadStdOut
// read stdout pipe until EOF.
//
func ReadStdOut(pipe io.ReadCloser) {

	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		log.Printf("stdout=%s", scanner.Text())
	}
}

// ReadStdErr
// read stderr pipe until EOF.
//
func ReadStdErr(pipe io.ReadCloser) {

	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		log.Printf("stderr=%s", scanner.Text())
	}
}

func RunCmdStdout(cmd *exec.Cmd) {

	stdOutPipe, err := cmd.StdoutPipe()
	CheckError(err)
	go ReadStdOut(stdOutPipe)

	stdErrPipe, err := cmd.StderrPipe()
	CheckError(err)
	go ReadStdErr(stdErrPipe)

	err = cmd.Run()
	CheckError(err)
}

func RunCmdStdoutIgnoreErr(cmd *exec.Cmd) {

	stdOutPipe, err := cmd.StdoutPipe()
	CheckError(err)
	go ReadStdOut(stdOutPipe)

	stdErrPipe, err := cmd.StderrPipe()
	CheckError(err)
	go ReadStdErr(stdErrPipe)

	_ = cmd.Run()
}
