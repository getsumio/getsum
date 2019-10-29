package supplier

import (
	"runtime"

	. "github.com/getsumio/getsum/internal/config"
	. "github.com/getsumio/getsum/internal/file"
	"github.com/getsumio/getsum/internal/status"
)

type ISupplierFactory interface {
	GetSupplierByAlgo(config *Config, algorithm *Algorithm) Supplier
}

type SupplierFactory struct {
}

func (factory *SupplierFactory) GetSupplierByAlgo(config *Config, algorithm *Algorithm) Supplier {

	return getSupplierInstance(config, algorithm)

}

func getSupplierInstance(config *Config, algo *Algorithm) Supplier {
	if *config.Supplier == "go" {
		s := &GoSupplier{}
		setFields(&s.BaseSupplier, *algo, config)
		return s
	} else if *config.Supplier == "openssl" {
		s := &CommandSupplier{Type: OPENSSL}
		setFields(&s.BaseSupplier, *algo, config)
		return s

	}
	switch runtime.GOOS {
	case "linux", "mac":
		s := &CommandSupplier{Type: UNIX}
		setFields(&s.BaseSupplier, *algo, config)
		return s
	case "windows":
		s := &CommandSupplier{Type: WINDOWS}
		setFields(&s.BaseSupplier, *algo, config)
		return s
	default:
		return nil

	}
}

func setFields(base *BaseSupplier, algo Algorithm, config *Config) {
	stat := &status.Status{status.PREPARED, "", ""}
	base.Algorithm = algo
	base.Key = *config.Key
	base.File = &File{Url: *config.File, Status: stat, Proxy: *config.Proxy, StoragePath: *config.Dir}
	base.TimeOut = *config.Timeout
	base.status = stat

}
