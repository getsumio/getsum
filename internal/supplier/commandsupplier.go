package supplier

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/getsumio/getsum/internal/status"
)

type CommandType int

const (
	UNIX CommandType = iota
	MAC
	WINDOWS
	OPENSSL
)

var unixAlgos []Algorithm = []Algorithm{MD5, SHA1, SHA224, SHA256, SHA384, SHA512}
var winAlgos []Algorithm = []Algorithm{MD2, MD4, MD5, SHA1, SHA224, SHA256, SHA384, SHA512}

var openSSLAlgos []Algorithm = []Algorithm{
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

type CommandSupplier struct {
	BaseSupplier
	Type CommandType
}

func (s *CommandSupplier) Supports() []Algorithm {
	switch s.Type {
	case OPENSSL:
		return openSSLAlgos
	case WINDOWS:
		return winAlgos
	default:
		return unixAlgos
	}
}

func execute(cmd *exec.Cmd, status chan string, e chan string) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		e <- err.Error()
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

func (s *CommandSupplier) Run(deleteOnExit bool) {
	if deleteOnExit {
		s.File.Delete()
	}
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
	e := make(chan string)
	cmd = getCommand(s)
	go execute(cmd, stat, e)
	for {
		tEnd := time.Now()
		took := tEnd.Sub(tStart)

		select {
		case <-t:
			kill(cmd)
			s.status.Type = status.TIMEDOUT
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			return
		case val := <-stat:
			s.status.Type = status.COMPLETED
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			if s.Type == OPENSSL {
				s.status.Checksum = strings.Fields(val)[1]
			} else {
				s.status.Checksum = strings.Fields(val)[0]
			}
			return
		case val := <-e:
			s.status.Type = status.ERROR
			s.status.Value = val
			return
		default:
			s.status.Type = status.RUNNING
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			time.Sleep(15 * time.Millisecond)
		}
	}
}

func (s *CommandSupplier) Status() *status.Status {
	return s.status
}

func (s *CommandSupplier) Terminate() {
	kill(cmd)
	if s.status.Type == status.RUNNING {
		s.status.Type = status.TERMINATED
	}
}

func getSSLCommand(algo Algorithm, path string) *exec.Cmd {
	algorithm := strings.ToLower(algo.Name())
	param := []string{"-", algorithm}
	algoParam := strings.Join(param, "")
	return exec.Command("openssl", "dgst", algoParam, path)
}

func getUnixCommand(a Algorithm, path string) *exec.Cmd {
	algo := strings.ToLower(a.Name())
	strs := []string{algo, "sum"}
	cmd := strings.Join(strs, "")
	return exec.Command(cmd, path)

}

func getWinCommand(a Algorithm, path string) *exec.Cmd {
	algo := strings.ToUpper(a.Name())
	return exec.Command("certUtil", "-hashfile", path, algo)

}

func getCommand(s *CommandSupplier) *exec.Cmd {
	switch s.Type {
	case WINDOWS:
		return getWinCommand(s.Algorithm, s.File.Path())
	case UNIX, MAC:
		return getUnixCommand(s.Algorithm, s.File.Path())
	case OPENSSL:
		return getSSLCommand(s.Algorithm, s.File.Path())
	default:
		panic("Unsupported command type!")
	}

}
