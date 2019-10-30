package providers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	. "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider/types"
	"github.com/getsumio/getsum/internal/status"
	. "github.com/getsumio/getsum/internal/supplier"
)

//Factory interface to init providers
type IProviderFactory interface {
	GetProviders(config *Config) (*Providers, error)
}

//provider factory
type ProviderFactory struct {
}

//Reads configuration file
//according to params initializes list of providers
//if present remote ones are AWS,GCE,IBM,ORCL,Azure
//local provider is Local
func (p *ProviderFactory) GetProviders(config *Config) (*Providers, error) {
	logger.Debug("Providers requested for config %v", *config)
	var factory ISupplierFactory = new(SupplierFactory)
	localProviders, err := getLocalProviders(config, factory)
	logger.Debug("Localproviders %v", localProviders)
	if err != nil {
		return nil, err
	}
	remoteProviders := getRemoteProviders(config)
	logger.Debug("Remoteproviders %v", remoteProviders)
	allProviders := append(localProviders, remoteProviders...)
	lengthLocal := len(localProviders)
	lengthRemote := len(remoteProviders)
	lengthTotal := len(allProviders)
	providers := &Providers{
		Locals:    localProviders,
		Remotes:   remoteProviders,
		All:       allProviders,
		HasRemote: lengthRemote > 0,
		HasLocal:  lengthLocal > 0,
		Length:    lengthTotal,
		Statuses:  make([]*status.Status, lengthTotal),
	}
	providers.HasValidation = *config.Cheksum != ""
	return providers, nil
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
func getRemoteProviders(config *Config) []*Provider {
	list := []*Provider{}
	if !*config.LocalOnly {
		for _, s := range config.Servers.Servers {
			provider := getRemoteProvider(config, &s)
			list = append(list, &provider)
		}
	}
	return list
}

func getRemoteProvider(config *Config, serverConfig *ServerConfig) Provider {
	r := &RemoteProvider{}
	r.Name = strings.Join([]string{serverConfig.Name, config.Algorithm[0]}, "-")
	r.address = serverConfig.Address
	r.client = getHttpClient(config)
	r.config = config
	r.Type = OnPremise
	r.WG = &sync.WaitGroup{}
	return r

}

func getProvider(pType ProviderType, supplier Supplier, config *Config, a Algorithm) Provider {
	l := &LocalProvider{}
	l.Name = strings.Join([]string{pType.Name(), a.Name()}, "-")
	l.Type = pType
	l.Supplier = supplier
	l.WG = &sync.WaitGroup{}
	return l

}

func getLocalProviders(config *Config, factory ISupplierFactory) ([]*Provider, error) {
	logger.Debug("Instantiating local providers")
	locals := []*Provider{}
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
			locals = append(locals, &l)
		}
	}

	return locals, nil
}
