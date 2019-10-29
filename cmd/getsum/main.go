// Package main provides ...
package main

import (
	"os"
	"os/signal"
	"time"

	parser "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider"
	"github.com/getsumio/getsum/internal/servers"
	"github.com/getsumio/getsum/internal/status"
	validator "github.com/getsumio/getsum/internal/validation"
)

func main() {
	config, err := parser.ParseConfig()
	if err != nil {
		logger.Error("Can not parse configuration: %s", err.Error())
		os.Exit(1)
	}

	validator.ValidateConfig(config)
	logger.SetLevel(*config.LogLevel)
	logger.Debug("Application  started, using config %v", *config)
	if *config.Serve {
		logger.Warn("Running in server mode listen address %s , port: %d", *config.Listen, *config.Port)
		server := &servers.OnPremiseServer{}
		err := server.Start(config)
		if err != nil {
			logger.Error("Can not start server: %s", err.Error())
		}
		os.Exit(1)
	}

	logger.Trace("Collecting providers")
	var factory IProviderFactory = new(ProviderFactory)
	providers, err := factory.GetProviders(config)
	length := len(providers)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	} else if length < 1 {
		logger.Error("There is no other supported algorithm asked to run, terminating")
		os.Exit(1)
	}

	logger.Info("providers: %v with size %d", providers, length)

	quit, wait := make(chan bool), make(chan bool)
	handleExit(quit)

	chans := make([]<-chan *status.Status, length)
	logger.Debug("Running providers, total length: %d", length)
	for i := 0; i < length; i++ {
		chans[i] = providers[i].Run(quit, wait)
	}
	stats := make([]*status.Status, length)
	logger.Trace("Starting to watch running processes")
	anyRunner, hasValidation, hasError := true, *config.Cheksum != "", false
	logger.Header(providers)
	for anyRunner {
		anyRunner = false
		for i := 0; i < length; i++ {
			s := <-chans[i]
			logger.Trace("Update value %v from provider", *s)
			if s.Type < status.COMPLETED {
				anyRunner = true
			} else if s.Type > status.COMPLETED {
				hasError = true
			} else {
				if hasValidation && s.Checksum != *config.Cheksum {
					logger.Debug("Checksum mismatch: asked: %s, found %s", *config.Cheksum, s.Checksum)
					s.Type = status.MISMATCH
					hasError = true
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

	if hasError {
		os.Exit(1)
	}

}

func handleExit(quit chan bool) {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt)

	go func() {
		<-sign
		quit <- true
		time.Sleep(time.Second)
		logger.Warn("\n\nTerminate requested by user")
		os.Exit(1)
	}()

}
