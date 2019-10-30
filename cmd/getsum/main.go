// Main package to handle configuration, parsing, validation and calculate checksums
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
	//read command line params as well as parse user config if exist
	config, err := parser.ParseConfig()
	if err != nil {
		logger.Error("Can not parse configuration: %s", err.Error())
		os.Exit(1)
	}

	//if all good validate configuration
	err = validator.ValidateConfig(config, false)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//if user defined set log level
	logger.SetLevel(*config.LogLevel)
	logger.Debug("Application  started, using config %v", *config)

	//check if user wants to run in listen mode
	if *config.Serve {
		logger.Warn("Running in server mode listen address %s , port: %d", *config.Listen, *config.Port)
		//onpremise is default server reads the config and listens given address and port no interface support atm
		server := &servers.OnPremiseServer{StoragePath: *config.Dir}
		err := server.Start(config)
		if err != nil {
			logger.Error("Can not start server: %s", err.Error())
		}
		os.Exit(1)
	}

	//providers are listeners of running processes
	//currently we have os,openssl and go
	//in future there might be aws,get etc... cloud function listeners
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

	//print the header
	logger.Header(providers)

	//if user specified a checksum to validate run in validation mode
	if providers.HasValidation && providers.HasRemote && providers.HasLocal {
		runRemoteFirst(providers, config)
	} else {
		runAll(providers, config)
	}
	logger.Debug("Application finish")

	//if error occured notify user
	if providers.HasError() {
		os.Exit(1)
	}

}

//simply loops until no running process found
//definition of running is routine status set to >= COMPLETED
func watch(providers *Providers, config *parser.Config) {
	for providers.IsRunning() {
		logger.Status(providers)
		//TODO: In case of remote http request set every 200ms
		//This might be needed to increase maybe to lower net load
		time.Sleep(200 * time.Millisecond)
	}
	//one final time print statuses
	logger.Status(providers)
}

//if user provided a validation checksum
// 1- check if any process completed
// 2- check if calculated sum matches
// 3- check if its matches if yes set status VALIDATED otherwise MISMATCH
func checkMismatch(providers *Providers, config *parser.Config) {
	if providers.HasMismatch(*config.Cheksum) {
		logger.Debug("There are mismatches")
		logger.Status(providers)
		//print results
		logger.Logsum(providers)
		//check if user still wants to keep the file
		if !*config.Keep {
			logger.Debug("Exiting application no keep setted")
			providers.Delete()
			os.Exit(1)
		}
	} else {
		//all good print one last time statuses
		logger.Status(providers)
	}

}

//this function runs all local/remote routines together
//disregarding user wants validation or not
//which means file still downloaded locally
//yet still if user dont set -keep param
//file will be deleted in case of validation present
func runAll(providers *Providers, config *parser.Config) {
	logger.Debug("Running all providers validation")
	providers.Run()
	watch(providers, config)
	//all process finish notify runners
	//to terminate their processes
	providers.Terminate(true)
	checkMismatch(providers, config)
	//print results
	logger.Logsum(providers)

}

//First runs remote runners
//which will reach servers and collect checksums
//if mismatch prevents host pc to fetch the file
//otherwise local will run after remote validation
//host will download the file and if mismatch
//will keep the file according to -keep param
func runRemoteFirst(providers *Providers, config *parser.Config) {
	logger.Debug("Running remote providers")
	//local waitgroup will be in wait mode
	providers.SuspendLocals()
	providers.Run()
	watch(providers, config)
	//terminates only finished processes
	providers.Terminate(false)
	checkMismatch(providers, config)
	//if there is mismatch os.Exit(1)
	//so no need to worry about suspended runners
	logger.Debug("Running local providers")
	providers.ResumeLocals()
	watch(providers, config)
	providers.Terminate(true)
	checkMismatch(providers, config)
	logger.Logsum(providers)
}

//in case of user in hurry or terminated
//make sure no runners
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
