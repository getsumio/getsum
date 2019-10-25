// Package main provides ...
package main

import (
	"os"
	"os/signal"
	"time"

	parser "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider"
	"github.com/getsumio/getsum/internal/status"
	validator "github.com/getsumio/getsum/internal/validation"
)

func main() {
	config := parser.ParseConfig()
	validator.ValidateConfig(config)
	logger.SetLevel(*config.LogLevel)
	logger.Debug("Application  started, using config %v", *config)
	logger.Trace("Collecting providers")
	var factory IProviderFactory = new(ProviderFactory)
	providers, err := factory.GetProviders(config)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Debug("providers: %v", providers)

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt)

	quit := make(chan bool)
	wait := make(chan bool)
	length := len(providers)
	chans := make([]<-chan *status.Status, length)
	logger.Debug("Running providers, total length: %d", length)
	for i := 0; i < length; i++ {
		chans[i] = providers[i].Run(quit, wait)
	}
	var anyRunner bool = true
	stats := make([]*status.Status, length)
	logger.Trace("Starting to watch running processes")
	hasValidation := *config.Cheksum != ""
	go func() {
		<-sign
		quit <- true
		time.Sleep(time.Second)
		logger.Warn("\n\nTerminate requested by user")
		os.Exit(1)
	}()
	logger.Header(providers)
	hasMisMatch := false
	for anyRunner {
		anyRunner = false
		for i := 0; i < length; i++ {
			s := <-chans[i]
			logger.Trace("Update value %v from provider", *s)
			if s.Type == status.PREPARED || s.Type == status.RUNNING || s.Type == status.STARTED || s.Type == status.DOWNLOAD {
				anyRunner = true
			} else if hasValidation && s.Type == status.COMPLETED {
				if s.Checksum != *config.Cheksum {
					hasMisMatch = true
					s.Type = status.MISMATCH
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
	if hasMisMatch {
		os.Exit(1)
	}

}
