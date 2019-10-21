package providers

import (
	"fmt"

	. "github.com/getsumio/getsum/internal/algorithm/supplier"
	. "github.com/getsumio/getsum/internal/config"
)

type IProviderFactory interface {
	GetProviders(config *Config) []Provider
}

type ProviderFactory struct {
}

func (p *ProviderFactory) GetProviders(config *Config) []Provider {
	factory := new(SupplierFactory)
	supplier := factory.GetSupplier(config)
	list := []Provider{}
	for i := 0; i < 6; i++ {
		l := &LocalProvider{}
		l.Supplier = supplier
		l.Name = fmt.Sprintf("local-pc%d", i)
		l.Proxy = config.Proxy
		l.File = config.File
		l.Type = Local
		list = append(list, l)
	}

	return list
}
