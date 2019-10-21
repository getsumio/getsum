// Package main provides ...
package main

import (
	. "github.com/getsumio/getsum/internal/algorithm/supplier"
	parser "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/providers"
	validator "github.com/getsumio/getsum/internal/validation"
)

func main() {
	config := parser.ParseConfig()
	validator.ValidateConfig(config)
	var factory IProviderFactory = new(ProviderFactory)
	var providers []Provider = factory.GetProviders(config)

	logger.Header(providers)

	quit := make(chan bool)
	wait := make(chan bool)
	length := len(providers)
	chans := make([]<-chan *Status, length)
	for i := 0; i < length; i++ {
		chans[i] = providers[i].Run(quit, wait)
	}
	var anyRunner bool = true
	stats := make([]*Status, length)
	for anyRunner {
		anyRunner = false
		for i := 0; i < length; i++ {
			s := <-chans[i]
			if s.Status == "RUNNING" || s.Status == "STARTED" {
				anyRunner = true
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

}
