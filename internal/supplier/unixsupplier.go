package supplier

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/getsumio/getsum/internal/status"
)

type UnixSupplier struct {
	BaseSupplier
}

func (s *UnixSupplier) Supports() []Algorithm {
	return []Algorithm{MD5, SHA1, SHA224, SHA256, SHA384, SHA512}
}

func execute(cmd *exec.Cmd, status chan string) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		status <- err.Error()

	} else {
		status <- string(out)
	}
}

func kill(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
	}
}

var cmd *exec.Cmd

func (s *UnixSupplier) Run() {

	err := s.File.Fetch(s.TimeOut)
	if err != nil {
		s.status.Value = err.Error()
		s.status.Type = status.ERROR
		return
	}

	tStart := time.Now()
	s.status.Type = status.STARTED
	t := time.After(time.Duration(s.TimeOut) * time.Second)
	stat := make(chan string)
	cmd = getCommand(s)
	go execute(cmd, stat)
	for {
		select {
		case <-t:
			kill(cmd)
			tEnd := time.Now()
			took := tEnd.Sub(tStart)
			s.status.Type = status.TIMEDOUT
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			return
		case val := <-stat:
			tEnd := time.Now()
			took := tEnd.Sub(tStart)
			s.status.Type = status.COMPLETED
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			s.status.Checksum = strings.Fields(val)[0]
			return
		default:
			tEnd := time.Now()
			took := tEnd.Sub(tStart)
			s.status.Type = status.RUNNING
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			time.Sleep(15 * time.Millisecond)
		}
	}
}

func (s *UnixSupplier) Status() *status.Status {
	return s.status
}

func (s *UnixSupplier) Terminate() {
	kill(cmd)
	if s.status.Type == status.RUNNING {
		s.status.Type = status.TERMINATED
	}
}

func getCommand(s *UnixSupplier) *exec.Cmd {
	algo := strings.ToLower(s.Algorithm.Name())
	strs := []string{algo, "sum"}
	cmd := strings.Join(strs, "")
	return exec.Command(cmd, s.File.Path())

}
