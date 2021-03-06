package providers

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

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
//if present remote ones are OnPremise,AWS,GCE,IBM,ORCL,Azure
//local provider is Local
func (p *ProviderFactory) GetProviders(config *Config) (*Providers, error) {
	logger.Debug("Providers requested for config %v", *config)
	//each provider (runner) has a supplier
	//providers just manages suppliers
	var factory ISupplierFactory = new(SupplierFactory)
	//collect local providers
	//there might be more than one i.e. user wants to run
	//all algos
	localProviders, err := getLocalProviders(config, factory)
	logger.Debug("Localproviders %v", localProviders)
	if err != nil {
		return nil, err
	}
	//collect remote providers
	remoteProviders := getRemoteProviders(config)
	logger.Debug("Remoteproviders %v", remoteProviders)
	//Providers struct is wrapper for both local remote
	//collect and set its fields
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
	providers.Filename = config.File
	return providers, nil
}

//utility to instantiate http client
//wrapped with proxy and timeout setting
func getHttpClient(config *Config) *http.Client {
	proxyUrl := http.ProxyFromEnvironment
	if *config.Proxy != "" {
		proxy, _ := url.Parse(*config.Proxy)
		proxyUrl = http.ProxyURL(proxy)
	}
	tr := &http.Transport{
		Proxy:           proxyUrl,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *config.InsecureSkipVerify},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(*config.Timeout) * time.Second,
	}
	return client
}

//if config is not local only and some servers present
//creates provider per server
func getRemoteProviders(config *Config) []*Provider {
	list := []*Provider{}
	//check user wants to run only on host pc
	if !*config.LocalOnly {
		for _, s := range config.Servers.Servers {
			provider := getRemoteProvider(config, &s)
			list = append(list, &provider)
		}
	}
	return list
}

//utility to build provider instance
//only OnPremise supported as of now
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

//utility to instantiate local provider
func getProvider(pType ProviderType, supplier Supplier, config *Config, a Algorithm) Provider {
	l := &LocalProvider{}
	l.Name = strings.Join([]string{pType.Name(), a.Name()}, "-")
	l.Type = pType
	l.Supplier = supplier
	l.WG = &sync.WaitGroup{}
	return l

}

//creates LocalProvider per algorithm specified in config
//in case of remoteOnly returns empty
//i.e. for param -a MD5, SHA512
//there will be 2 provider with MD5 supplier and SHA512 supplier
//supplier factory will take care of what algo and which library will be used
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
			supplier, err := factory.GetSupplierByAlgo(config, &a)
			if err != nil {
				logger.Warn(err.Error())
				continue
			}

			l := getProvider(Local, supplier, config, a)
			logger.Debug("Generated provider: %v", l)
			locals = append(locals, &l)
		}
	}

	return locals, nil
}
