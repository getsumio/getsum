package supplier

import (
	. "github.com/getsumio/getsum/internal/config"
	. "github.com/getsumio/getsum/internal/file"
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

	stat := &Status{"PREPARED", "", ""}
	s := &UnixSupplier{}
	s.Algorithm = *algorithm
	s.File = &File{Url: *config.File, Status: stat}
	s.TimeOut = *config.Timeout
	s.status = stat
	return s

}
