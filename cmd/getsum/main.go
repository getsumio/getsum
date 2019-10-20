// Package main provides ...
package main

import (
	"fmt"

	. "github.com/getsumio/getsum/internal/algorithm"
	. "github.com/getsumio/getsum/internal/algorithm/supplier"
	parser "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/providers"
	validator "github.com/getsumio/getsum/internal/validation"
)

func main() {
	config := parser.ParseConfig()
	validator.ValidateConfig(config)
	list := []Provider{}
	for i := 0; i < 6; i++ {
		l := new(LocalProvider)
		s := new(UnixSupplier)
		s.Algorithm = SHA3
		l.Supplier = s
		l.Name = fmt.Sprintf("local-pc%d", i)
		list = append(list, l)
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
