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
	stat := &status.Status{status.PREPARED, "", ""}
	if *config.Supplier == "go" {
		s := &GoSupplier{}
		s.Algorithm = *algo
		s.Key = *config.Key
		s.File = &File{Url: *config.File, Status: stat}
		s.TimeOut = *config.Timeout
		s.status = stat
		return s
	}
	switch runtime.GOOS {
	case "linux":
		s := &UnixSupplier{}
		s.Algorithm = *algo
		s.Key = *config.Key
		s.File = &File{Url: *config.File, Status: stat}
		s.TimeOut = *config.Timeout
		s.status = stat
		return s
	default:
		return nil

	}
}

func setFields(base *BaseSupplier, algo Algorithm, config *Config) {
	stat := &status.Status{status.PREPARED, "", ""}
	base.Algorithm = algo
	base.Key = *config.Key
	base.File = &File{Url: *config.File, Status: stat}
	base.TimeOut = *config.Timeout
	base.status = stat

}
