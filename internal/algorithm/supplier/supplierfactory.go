package supplier

import (
	. "github.com/getsumio/getsum/internal/config"
)

type ISupplierFactory interface {
	GetSupplier(config *Config) Supplier
}

type SupplierFactory struct {
}

func (factory *SupplierFactory) GetSupplier(config *Config) Supplier {

	s := &UnixSupplier{}
	s.Algorithm = *config.Algorithm
	s.File = *config.File
	s.TimeOut = *config.Timeout
	s.status = &Status{"PREPARED", "", ""}
	return s

}
