package supplier

import (
	. "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/file"
)

type ISupplierFactory interface {
	GetSupplier(config *Config) Supplier
	GetSupplierByAlgo(config *Config, algorithm *Algorithm) Supplier
}

type SupplierFactory struct {
}

func (factory *SupplierFactory) GetSupplier(config *Config) Supplier {

	algorithm, _ := ValueOf(config.Algorithm)

	return factory.GetSupplierByAlgo(config, &algorithm)

}

func (factory *SupplierFactory) GetSupplierByAlgo(config *Config, algorithm *Algorithm) Supplier {

	s := &UnixSupplier{}
	s.Algorithm = *algorithm
	s.File = &file.File{Url: *config.File}
	s.TimeOut = *config.Timeout
	s.status = &Status{"PREPARED", "", ""}
	return s

}
