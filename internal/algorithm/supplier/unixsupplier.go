package supplier

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type UnixSupplier struct {
	BaseSupplier
}

func execute(cmd *exec.Cmd, quit <-chan bool, status chan string) {
	go func() {
		for {
			select {
			case <-quit:
				if cmd.Process != nil {
					cmd.Process.Kill()
				}
			}
		}
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		status <- err.Error()

	} else {
		status <- string(out)
	}
}

var quit chan bool

func (s *UnixSupplier) Run() {

	sVal := "DOWNLOADING"
	s.File.Status = sVal
	s.status.Status = sVal
	s.status.Value = "0%"
	go func() {
		for s.File.Status == sVal {
			s.status.Value = s.File.StatusValue
			time.Sleep(time.Second)
		}

	}()
	err := s.File.Fetch(s.TimeOut)
	if err != nil {
		s.status.Status = "ERROR"
		s.status.Value = err.Error()
		return
	}

	tStart := time.Now()
	s.status.Status = "STARTED"
	t := time.After(time.Duration(s.TimeOut) * time.Second)
	quit = make(chan bool)
	status := make(chan string)
	go execute(getCommand(s), quit, status)
	for {
		select {
		case <-t:
			tEnd := time.Now()
			took := tEnd.Sub(tStart)
			s.status.Status = "TIMEDOUT"
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			quit <- true
			return
		case val := <-status:
			tEnd := time.Now()
			took := tEnd.Sub(tStart)
			s.status.Status = "COMPLETED"
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			s.status.Checksum = val
			return
		default:
			tEnd := time.Now()
			took := tEnd.Sub(tStart)
			s.status.Status = "RUNNING"
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			time.Sleep(15 * time.Millisecond)
		}
	}
}

func (s *UnixSupplier) Status() *Status {
	return s.status
}

func (s *UnixSupplier) Terminate() {
	s.status = &Status{"Terminated", "", ""}
	quit <- true
}

func getCommand(s *UnixSupplier) *exec.Cmd {
	algo := strings.ToLower(s.Algorithm)
	strs := []string{algo, "sum"}
	cmd := strings.Join(strs, "")
	return exec.Command(cmd, s.File.Path())

}
