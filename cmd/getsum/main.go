// Package main provides ...
package main

import (
	"fmt"

	. "github.com/getsumio/getsum/internal/algorithm"
	. "github.com/getsumio/getsum/internal/algorithm/supplier"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/providers"
)

func main() {
	logger.Debug("test")
	logger.Level = logger.LevelError
	logger.Debug("test2")
	logger.Level = logger.LevelTrace
	logger.Trace("Application started")
	logger.Debug("Fetching configuration file")
	logger.Info("Configuration content: asdsad")
	logger.Warn("There are unrecognized settings on config ignoring")
	logger.Error("Config supplier is required!!!")

	list := []Provider{}
	for i := 0; i < 5; i++ {
		l := new(LocalProvider)
		s := new(UnixSupplier)
		s.Algorithm = SHA3
		l.Supplier = s
		l.Name = fmt.Sprintf("local-pc%d", i)
		list = append(list, l)
	}

	for _, i := range list {

		logger.Info(i.Data().Name)
	}

	logger.Header(list)

	quit := make(chan bool)
	wait := make(chan bool)
	chans := []<-chan *Status{}
	for _, p := range list {
		chans = append(chans, p.Run(quit, wait))
	}
	var anyRunner bool = true
	for anyRunner {
		anyRunner = false
		stats := []*Status{}
		for _, c := range chans {
			if c == nil {
				logger.Error("channel is nil!")
			}
			s := <-c
			if s.Status == "RUNNING" || s.Status == "STARTED" {
				anyRunner = true
			}

			stats = append(stats, s)

		}
		if anyRunner {
			wait <- true
		} else {
			quit <- true
		}
		logger.Status(stats)
	}

}
