package providers

import (
	"fmt"

	. "github.com/getsumio/getsum/internal/algorithm/supplier"
	. "github.com/getsumio/getsum/internal/config"
)

//Factory interface to init providers
type IProviderFactory interface {
	GetProviders(config *Config) []Provider
}

//provider factory
type ProviderFactory struct {
}

//Reads configuration file
//according to params initializes list of providers
//if present remote ones are AWS,GCE,IBM,ORCL,Azure
//local provider is Local
func (p *ProviderFactory) GetProviders(config *Config) []Provider {
	var factory ISupplierFactory = new(SupplierFactory)
	list := []Provider{}
	list = append(list, getLocalProviders(config, factory)...)

	return list
}

func getLocalProviders(config *Config, factory ISupplierFactory) []Provider {
	locals := []Provider{}
	if !*config.RemoteOnly {
		if *config.All {
			for _, a := range Algorithms {
				l := &LocalProvider{}
				l.Supplier = factory.GetSupplierByAlgo(config, &a)
				l.Name = fmt.Sprintf("local-%s", a.Name())
				l.Proxy = config.Proxy
				l.File = config.File
				l.Type = Local
				locals = append(locals, l)
			}

		} else {
			l := &LocalProvider{}
			l.Supplier = factory.GetSupplier(config)
			l.Name = "local-pc"
			l.Proxy = config.Proxy
			l.File = config.File
			l.Type = Local
			locals = append(locals, l)

		}
	}

	return locals
}
