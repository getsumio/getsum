// Package main provides ...
package main

import (
	"fmt"
	"time"

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

	for i := 0; i < 10; i++ {
		for _, p := range list {
			p.Data().Value = fmt.Sprintf("%d%%", i*10)
			p.Run()
			//logger.Inplace("Processing %d %s", i, sum)
		}
		logger.Status(list)
		time.Sleep(time.Second)
	}

}
