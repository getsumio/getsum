package supplier

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/getsumio/getsum/internal/file"
	"github.com/getsumio/getsum/internal/status"
)

type GoSupplier struct {
	BaseSupplier
}

func (s *GoSupplier) Run() {
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
	go calculate(s.Algorithm, stat, s.File)
	for {
		select {
		case <-t:
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

func (s *GoSupplier) Status() *status.Status {
	return s.status
}

func (s *GoSupplier) Terminate() {
	if s.status.Type == status.RUNNING {
		s.status.Type = status.TERMINATED
	}

}

func calculate(algo Algorithm, status chan string, file *file.File) {
	sha256.New()

}
