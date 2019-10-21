package providers

import (
	. "github.com/getsumio/getsum/internal/algorithm/supplier"
)

type Provider interface {
	Close()
	Data() *BaseProvider
	Run(quit <-chan bool, wait <-chan bool) <-chan *Status
}

type BaseProvider struct {
	Name     string
	Address  *string
	Zone     *string
	Proxy    *string
	Type     ProviderType
	Supplier Supplier
}

type ProviderType int

const (
	Aws = iota
	Google
	Oracle
	Azure
	IBM
	Local
)
