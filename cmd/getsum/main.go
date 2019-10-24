// Package main provides ...
package main

import (
	parser "github.com/getsumio/getsum/internal/config"
	. "github.com/getsumio/getsum/internal/file"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/providers"
	validator "github.com/getsumio/getsum/internal/validation"
)

func main() {
	config := parser.ParseConfig()
	validator.ValidateConfig(config)
	logger.SetLevel(*config.LogLevel)
	logger.Debug("Application  started, using config %v", *config)
	logger.Trace("Collecting providers")
	var factory IProviderFactory = new(ProviderFactory)
	var providers []Provider = factory.GetProviders(config)
	logger.Debug("providers: %v", providers)

	logger.Header(providers)

	quit := make(chan bool)
	wait := make(chan bool)
	length := len(providers)
	chans := make([]<-chan *Status, length)
	logger.Debug("Running providers, total length: %d", length)
	for i := 0; i < length; i++ {
		chans[i] = providers[i].Run(quit, wait)
	}
	var anyRunner bool = true
	stats := make([]*Status, length)
	logger.Trace("Starting to watch running processes")
	hasValidation := *config.Cheksum != ""
	for anyRunner {
		anyRunner = false
		for i := 0; i < length; i++ {
			s := <-chans[i]
			logger.Debug("Update value %v from provider", *s)
			if s.Status == "PREPARED" || s.Status == "RUNNING" || s.Status == "STARTED" || s.Status == "DOWNLOAD" {
				anyRunner = true
			} else if hasValidation && s.Status == "COMPLETED" {
				if s.Checksum != *config.Cheksum {
					s.Status = "MISMATCH"
				}

			}

			stats[i] = s

		}
		if anyRunner {
			wait <- true
		} else {
			quit <- true
		}
		logger.Status(stats)
	}
	logger.Logsum(providers, stats)
	logger.Debug("Application finish")

}
