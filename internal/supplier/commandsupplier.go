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
var winAlgos []Algorithm = []Algorithm{MD2, MD4, MD5, SHA1, SHA256, SHA384, SHA512}

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

//Common command executor for openssl or OS apps
//for unix, max i.e. shaXXXsum application will be called
//for openssl it is openssl see getCommand method for commands
type CommandSupplier struct {
	BaseSupplier
	Type CommandType
}

//check if this runner can run given algo
//i.e. openssl doesnt have MD2 support
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

//execute command on operating system
//if error occured update status
//then termination is listener job
func execute(cmd *exec.Cmd, status chan string, e chan string) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		e <- err.Error()
	} else {
		status <- string(out)
	}
}

//terminate called
//kill the running process
//if there is any
func kill(cmd *exec.Cmd) error {
	if cmd != nil && cmd.Process != nil {
		return cmd.Process.Kill()
	}
	return nil
}

var cmd *exec.Cmd

//start calculation
func (s *CommandSupplier) Run() {
	//fetch the file
	err := s.File.Fetch(s.TimeOut)
	if err != nil {
		s.status.Value = err.Error()
		s.status.Type = status.ERROR
		return
	}

	//make sure timeout not reached
	tStart := time.Now()
	s.status.Type = status.STARTED
	t := time.After(time.Duration(s.TimeOut) * time.Second)
	stat := make(chan string)
	e := make(chan string)
	cmd = getCommand(s)
	//execute command on routine in the main time start watching
	go execute(cmd, stat, e)
	for {
		tEnd := time.Now()
		took := tEnd.Sub(tStart)

		select {
		case <-t:
			//timeout occured kill and terminate
			kill(cmd)
			s.status.Type = status.TIMEDOUT
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			return
		case val := <-stat:
			//process completed and we receive a result
			//return result
			s.status.Type = status.COMPLETED
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			if s.Type == OPENSSL {
				s.status.Checksum = strings.Fields(val)[1]
			} else if s.Type == WINDOWS {
				s.status.Checksum = strings.Split(val, "\n")[1]
			} else {
				s.status.Checksum = strings.Fields(val)[0]
			}
			return
		case val := <-e:
			//we got error
			s.status.Type = status.ERROR
			s.status.Value = val
			return
		default:
			//still running update time
			s.status.Type = status.RUNNING
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			time.Sleep(15 * time.Millisecond)
		}
	}
}

//return status
func (s *CommandSupplier) Status() *status.Status {
	return s.status
}

//remove file
func (s *CommandSupplier) Delete() {
	s.File.Delete()
}

//terminate process
func (s *CommandSupplier) Terminate() error {
	err := kill(cmd)
	if s.status.Type == status.RUNNING {
		s.status.Type = status.TERMINATED
	}
	return err
}

//returns i.e. openssl dgst -sha512 /file/path
func getSSLCommand(algo Algorithm, path string) *exec.Cmd {
	algorithm := strings.ToLower(algo.Name())
	param := []string{"-", algorithm}
	algoParam := strings.Join(param, "")
	return exec.Command("openssl", "dgst", algoParam, path)
}

//returns i.e. sha512sum /file/path
func getUnixCommand(a Algorithm, path string) *exec.Cmd {
	algo := strings.ToLower(a.Name())
	strs := []string{algo, "sum"}
	cmd := strings.Join(strs, "")
	return exec.Command(cmd, path)

}

//returns i.e. certutil -hashfile /file/path SHA512
func getWinCommand(a Algorithm, path string) *exec.Cmd {
	algo := strings.ToUpper(a.Name())
	return exec.Command("certUtil", "-hashfile", path, algo)

}

//returns related command according to -lib selection by the user
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
