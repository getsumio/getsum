package providers

import (
	"fmt"
	"net/http"
	"net/url"

	. "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider/types"
	. "github.com/getsumio/getsum/internal/supplier"
)

//Factory interface to init providers
type IProviderFactory interface {
	GetProviders(config *Config) ([]Provider, error)
}

//provider factory
type ProviderFactory struct {
}

//Reads configuration file
//according to params initializes list of providers
//if present remote ones are AWS,GCE,IBM,ORCL,Azure
//local provider is Local
func (p *ProviderFactory) GetProviders(config *Config) ([]Provider, error) {
	logger.Debug("Providers requested for config %v", *config)
	var factory ISupplierFactory = new(SupplierFactory)
	list := []Provider{}
	localProviders, err := getLocalProviders(config, factory)
	if err != nil {
		return nil, err
	}

	list = append(list, localProviders...)
	list = append(list, getRemoteProviders(config)...)
	logger.Debug("Generated providers: %v", list)

	return list, nil
}

func getHttpClient(config *Config) *http.Client {
	proxyUrl := http.ProxyFromEnvironment
	if *config.Proxy != "" {
		proxy, _ := url.Parse(*config.Proxy)
		proxyUrl = http.ProxyURL(proxy)
	}
	tr := &http.Transport{
		Proxy: proxyUrl,
	}
	client := &http.Client{
		Transport: tr,
	}
	return client
}
func getRemoteProviders(config *Config) []Provider {
	list := []Provider{}
	for _, s := range config.Servers.Servers {
		list = append(list, getRemoteProvider(config, &s))
	}
	return list
}

func getRemoteProvider(config *Config, serverConfig *ServerConfig) Provider {
	r := &RemoteProvider{}
	r.Name = serverConfig.Name
	r.address = serverConfig.Address
	r.client = getHttpClient(config)
	r.config = config
	return r

}

func getProvider(pType ProviderType, supplier Supplier, config *Config, a Algorithm) Provider {
	l := &LocalProvider{}
	l.Name = fmt.Sprintf("%s-%s", pType.Name(), a.Name())
	l.Proxy = config.Proxy
	l.Type = pType
	l.Supplier = supplier
	return l

}

func getLocalProviders(config *Config, factory ISupplierFactory) ([]Provider, error) {
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
			supports := false
			for _, supportedAlgo := range supplier.Supports() {
				if supportedAlgo == a {
					supports = true
				}
			}
			if !supports {
				logger.Warn(fmt.Sprintf("Algorithm %s not supported for local provider using %s libraries", a.Name(), *config.Supplier))
				continue
			}
			l := getProvider(Local, supplier, config, a)
			logger.Debug("Generated provider: %v", l)
			locals = append(locals, l)
		}
	}

	return locals, nil
}
