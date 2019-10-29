// Package main provides ...
package main

import (
	"os"
	"os/signal"
	"time"

	parser "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider"
	. "github.com/getsumio/getsum/internal/provider/types"
	"github.com/getsumio/getsum/internal/servers"
	validator "github.com/getsumio/getsum/internal/validation"
)

func main() {
	config, err := parser.ParseConfig()
	if err != nil {
		logger.Error("Can not parse configuration: %s", err.Error())
		os.Exit(1)
	}

	err = validator.ValidateConfig(config, false)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

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
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	} else if providers.Length < 1 {
		logger.Error("No runner specified or no supported algorithm, terminating")
		os.Exit(1)
	}

	logger.Info("providers: %v with size %d", providers, providers.Length)
	handleExit(providers)

	logger.Trace("Starting to watch running processes")
	hasValidation := *config.Cheksum != ""

	logger.Header(providers)
	if hasValidation && providers.HasRemote && providers.HasLocal {
		runRemoteFirst(providers, config)
	} else {
		runAll(providers, config)
	}
	logger.Debug("Application finish")

	if providers.HasError() {
		os.Exit(1)
	}

}

func watch(providers *Providers, config *parser.Config) {
	for providers.IsRunning() {
		logger.Status(providers, *config.Cheksum)
		time.Sleep(200 * time.Millisecond)
	}
	logger.Status(providers, *config.Cheksum)
}

func checkMismatch(providers *Providers, config *parser.Config) {
	if providers.HasMismatch(*config.Cheksum) {
		logger.Debug("There are mismatches")
		logger.Status(providers, *config.Cheksum)
		logger.Logsum(providers.All, providers.Status())
		os.Exit(1)
	}

}

func runAll(providers *Providers, config *parser.Config) {
	logger.Debug("Running all providers validation")
	providers.Run()
	watch(providers, config)
	providers.Terminate(true)
	checkMismatch(providers, config)
	logger.Logsum(providers.All, providers.Status())

}

func runRemoteFirst(providers *Providers, config *parser.Config) {
	logger.Debug("Running remote providers")
	providers.SuspendLocals()
	providers.Run()
	watch(providers, config)
	providers.Terminate(false)
	checkMismatch(providers, config)
	logger.Debug("Running local providers")
	providers.ResumeLocals()
	watch(providers, config)
	providers.Terminate(true)
	checkMismatch(providers, config)
	logger.Status(providers, *config.Cheksum)
	logger.Logsum(providers.All, providers.Status())
}

func handleExit(providers *Providers) {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt)

	go func() {
		<-sign
		providers.Terminate(true)
		time.Sleep(300 * time.Millisecond)
		logger.Warn("\n\nTerminate requested by user")
		os.Exit(1)
	}()

}
