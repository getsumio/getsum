package supplier

import (
	"runtime"

	. "github.com/getsumio/getsum/internal/config"
	. "github.com/getsumio/getsum/internal/file"
	"github.com/getsumio/getsum/internal/status"
)

//reads the config and returns realted supplier
type ISupplierFactory interface {
	GetSupplierByAlgo(config *Config, algorithm *Algorithm) Supplier
}

//factory struct
type SupplierFactory struct {
}

//returns supplier instance for the given algo and lib
//i.e. for -lib go -a MD5 it will return GoSupplier to calculate MD5
func (factory *SupplierFactory) GetSupplierByAlgo(config *Config, algorithm *Algorithm) Supplier {

	return getSupplierInstance(config, algorithm)

}

var cache map[string]Supplier = make(map[string]Supplier)

//creates supplier instance
func getSupplierInstance(config *Config, algo *Algorithm) Supplier {
	if *config.Supplier == "go" {
		s, ok := cache["go"+string(*algo)]
		if !ok {
			s = &GoSupplier{}
			cache["go"+string(*algo)] = s

		}
		setFields(s.Data(), *algo, config)
		return s
	} else if *config.Supplier == "openssl" {
		s, ok := cache["openssl"+string(*algo)]
		if !ok {
			s = &CommandSupplier{Type: OPENSSL}
			cache["openssl"+string(*algo)] = s

		}
		setFields(s.Data(), *algo, config)
		return s
	}
	switch runtime.GOOS {
	case "linux", "mac":
		s, ok := cache["mac"+string(*algo)]
		if !ok {
			s = &CommandSupplier{Type: UNIX}
			cache["mac"+string(*algo)] = s
		}
		setFields(s.Data(), *algo, config)
		return s
	case "windows":
		s, ok := cache["windows"+string(*algo)]
		if !ok {
			s = &CommandSupplier{Type: WINDOWS}
			cache["windows"+string(*algo)] = s
		}
		setFields(s.Data(), *algo, config)
		return s
	default:
		return nil

	}
}

//utility to set fields
func setFields(base *BaseSupplier, algo Algorithm, config *Config) {
	if base.status == nil {
		base.status = &status.Status{}
	}
	base.status.Type = status.PREPARED
	base.status.Value = ""
	base.status.Checksum = ""
	base.Algorithm = algo
	base.Key = *config.Key
	if base.File == nil {
		base.File = &File{}
	}
	base.File.Reset()
	base.File.Url = *config.File
	base.File.Status = base.status
	base.File.Proxy = *config.Proxy
	base.File.StoragePath = *config.Dir

	base.TimeOut = *config.Timeout
}
