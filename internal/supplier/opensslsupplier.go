package supplier

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/getsumio/getsum/internal/status"
)

type OpenSSLSupplier struct {
	BaseSupplier
}

func (s *OpenSSLSupplier) Supports() []Algorithm {
	return []Algorithm{
		MD4,
		MD5,
		SHA1,
		SHA224,
		SHA256,
		SHA384,
		SHA512,
		RMD160,
		SHA3_224,
		SHA3_256,
		SHA3_384,
		SHA3_512,
		SHA512_224,
		SHA512_256,
		BLAKE2s256,
		BLAKE2b512,
		SHAKE128,
		SHAKE256,
		SM3,
	}

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

func (s *OpenSSLSupplier) Run() {

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
	cmd = getSSLCommand(s)
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
			s.status.Checksum = strings.Fields(val)[1]
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

func (s *OpenSSLSupplier) Status() *status.Status {
	return s.status
}

func (s *OpenSSLSupplier) Terminate() {
	kill(cmd)
	if s.status.Type == status.RUNNING {
		s.status.Type = status.TERMINATED
	}
}

func getSSLCommand(s *OpenSSLSupplier) *exec.Cmd {
	algo := strings.ToLower(s.Algorithm.Name())
	strs := []string{"-", algo}
	algoParam := strings.Join(strs, "")
	return exec.Command("openssl", "dgst", algoParam, s.File.Path())

}
