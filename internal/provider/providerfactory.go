package providers

import (
	"fmt"

	. "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider/types"
	. "github.com/getsumio/getsum/internal/supplier"
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
	logger.Debug("Providers requested for config %v", *config)
	var factory ISupplierFactory = new(SupplierFactory)
	list := []Provider{}
	list = append(list, getLocalProviders(config, factory)...)
	logger.Debug("Generated providers: %v", list)

	return list
}

func getProvider(pType ProviderType, supplier Supplier, config *Config, a Algorithm) Provider {
	l := &LocalProvider{}
	l.Name = fmt.Sprintf("%s %s-%s", l.Region(), pType.Name(), a.Name())
	l.Proxy = config.Proxy
	l.File = config.File
	l.Type = pType
	l.Supplier = supplier
	return l

}

func getLocalProviders(config *Config, factory ISupplierFactory) []Provider {
	logger.Debug("Instantiating local providers")
	locals := []Provider{}
	if !*config.RemoteOnly {
		logger.Debug("Config is remote only")
		var algos []Algorithm
		if *config.All {
			algos = Algorithms
		} else {
			for _, s := range config.Algorithm {
				algos = append(algos, ValueOf(&s))
			}
		}
		logger.Debug("User requests all algos to be runned")
		for _, a := range algos {
			logger.Debug("Creating local provider for algorithm %s", a.Name())
			supplier := factory.GetSupplierByAlgo(config, &a)
			l := getProvider(Local, supplier, config, a)
			logger.Debug("Generated provider: %v", l)
			locals = append(locals, l)
		}
	}

	return locals
}