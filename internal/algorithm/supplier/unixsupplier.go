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

func (s *UnixSupplier) Run() error {

	tStart := time.Now()
	stat := &Status{"STARTED", "", ""}
	s.status = stat
	t := time.After(time.Duration(s.TimeOut) * time.Second)
	quit = make(chan bool)
	status := make(chan string)
	go execute(getCommand(s), quit, status)
	for {
		select {
		case <-t:
			tEnd := time.Now()
			took := tEnd.Sub(tStart)
			stat.Status = "TIMEDOUT"
			stat.Value = fmt.Sprintf("%dms", took.Milliseconds())
			s.status = stat
			quit <- true
			return nil
		case val := <-status:
			tEnd := time.Now()
			took := tEnd.Sub(tStart)
			stat.Status = "COMPLETED"
			stat.Value = fmt.Sprintf("%dms", took.Milliseconds())
			stat.Checksum = val
			s.status = stat
			return nil
		default:
			tEnd := time.Now()
			took := tEnd.Sub(tStart)
			stat.Status = "RUNNING"
			stat.Value = fmt.Sprintf("%dms", took.Milliseconds())
			s.status = stat
			time.Sleep(15 * time.Millisecond)
		}
	}
	return nil
}

func (s *UnixSupplier) Status() *Status {
	return s.status
}

func (s *UnixSupplier) Terminate() error {
	s.status = &Status{"Terminated", "", ""}
	quit <- true
	return nil
}

func getCommand(s *UnixSupplier) *exec.Cmd {
	algo := strings.ToLower(s.Algorithm)
	strs := []string{algo, "sum"}
	cmd := strings.Join(strs, "")
	return exec.Command(cmd, s.File)

}
