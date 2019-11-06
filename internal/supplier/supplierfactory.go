package supplier

import (
	"errors"
	"fmt"
	"runtime"

	. "github.com/getsumio/getsum/internal/config"
	. "github.com/getsumio/getsum/internal/file"
	"github.com/getsumio/getsum/internal/status"
)

//reads the config and returns realted supplier
type ISupplierFactory interface {
	GetSupplierByAlgo(config *Config, algorithm *Algorithm) (Supplier, error)
}

//factory struct
type SupplierFactory struct {
}

//returns supplier instance for the given algo and lib
//i.e. for -lib go -a MD5 it will return GoSupplier to calculate MD5
func (factory *SupplierFactory) GetSupplierByAlgo(config *Config, algorithm *Algorithm) (Supplier, error) {

	supplier, err := getSupplierInstance(config, algorithm)
	if err != nil {
		return nil, err
	}
	if !isSupplierSupportsAlgorithm(supplier, algorithm) {
		msg := fmt.Sprintf("Algorithm %s not supported for provider using %s libraries", algorithm.Name(), *config.Supplier)
		return nil, errors.New(msg)
	}
	return supplier, nil

}

func isSupplierSupportsAlgorithm(supplier Supplier, algo *Algorithm) bool {
	supports := false
	for _, supportedAlgo := range supplier.Supports() {
		if supportedAlgo == *algo {
			supports = true
		}
	}
	return supports
}

//creates supplier instance
func getSupplierInstance(config *Config, algo *Algorithm) (Supplier, error) {
	if *config.Supplier == "go" {
		s := &GoSupplier{}
		setFields(s.Data(), *algo, config)
		return s, nil
	} else if *config.Supplier == "openssl" {
		s := &CommandSupplier{Type: OPENSSL}
		setFields(s.Data(), *algo, config)
		return s, nil
	}
	switch runtime.GOOS {
	case "linux", "mac":
		s := &CommandSupplier{Type: UNIX}
		setFields(s.Data(), *algo, config)
		return s, nil
	case "windows":
		s := &CommandSupplier{Type: WINDOWS}
		setFields(s.Data(), *algo, config)
		return s, nil
	default:
		msg := fmt.Sprintf("Unsupported library: %s", *config.Supplier)
		return nil, errors.New(msg)

	}
}

//utility to set fields
func setFields(base *BaseSupplier, algo Algorithm, config *Config) {
	base.status = &status.Status{}
	base.status.Type = status.PREPARED
	base.status.Value = ""
	base.status.Checksum = ""
	base.Algorithm = algo
	base.Key = *config.Key
	base.File = &File{}
	base.File.Reset()
	base.File.Url = *config.File
	base.File.Status = base.status
	base.File.Proxy = *config.Proxy
	base.File.StoragePath = *config.Dir
	base.File.SkipVerify = *config.InsecureSkipVerify

	base.TimeOut = *config.Timeout
	base.IsConcurrent = *config.All
}
