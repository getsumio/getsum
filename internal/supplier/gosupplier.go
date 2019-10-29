package supplier

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
	"time"

	"github.com/getsumio/getsum/internal/file"
	"github.com/getsumio/getsum/internal/status"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

type GoSupplier struct {
	BaseSupplier
}

func (s *GoSupplier) Supports() []Algorithm {
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
		BLAKE2b256,
		BLAKE2b384,
		BLAKE2b512,
	}
}

func (s *GoSupplier) Run(deleteOnExit bool) {
	if deleteOnExit {
		defer s.File.Delete()
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
	hash, err := getHash(s.Algorithm, s.Key)
	if err != nil {
		s.status.Value = err.Error()
		s.status.Type = status.ERROR
	}

	go calculate(hash, stat, s.File)
	for {
		tEnd := time.Now()
		took := tEnd.Sub(tStart)

		select {
		case <-t:
			s.status.Type = status.TIMEDOUT
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			return
		case val := <-stat:
			s.status.Type = status.COMPLETED
			s.status.Value = fmt.Sprintf("%dms", took.Milliseconds())
			s.status.Checksum = strings.Fields(val)[0]
			return
		default:
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

func calculate(hash hash.Hash, status chan string, file *file.File) {
	f, _ := os.Open(file.Path())
	defer f.Close()
	if _, err := io.Copy(hash, f); err != nil {
	}
	status <- hex.EncodeToString(hash.Sum(nil))
}

func getHash(algo Algorithm, key string) (hash.Hash, error) {
	switch algo {
	case MD4:
		return md4.New(), nil
	case MD5:
		return md5.New(), nil
	case SHA1:
		return sha1.New(), nil
	case SHA224:
		return sha256.New224(), nil
	case SHA256:
		return sha256.New(), nil
	case SHA384:
		return sha512.New384(), nil
	case SHA512:
		return sha512.New(), nil
	case RMD160:
		return ripemd160.New(), nil
	case SHA3_224:
		return sha3.New224(), nil
	case SHA3_256:
		return sha3.New256(), nil
	case SHA3_384:
		return sha3.New384(), nil
	case SHA3_512:
		return sha3.New512(), nil
	case SHA512_224:
		return sha512.New512_224(), nil
	case SHA512_256:
		return sha512.New512_256(), nil
	case BLAKE2s256:
		return blake2s.New256([]byte(key))
	case BLAKE2b256:
		return blake2b.New256([]byte(key))
	case BLAKE2b384:
		return blake2b.New384([]byte(key))
	case BLAKE2b512:
		return blake2b.New512([]byte(key))
	default:
		return nil, errors.New(fmt.Sprintf("Algorithm %s not supported", algo.Name()))
	}
}
