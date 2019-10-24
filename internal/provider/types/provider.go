package providers

import (
	"github.com/getsumio/getsum/internal/status"
	"github.com/getsumio/getsum/internal/supplier"
)

type Provider interface {
	Close()
	Data() *BaseProvider
	Run(quit <-chan bool, wait <-chan bool) <-chan *status.Status
}

type BaseProvider struct {
	Name     string
	Address  *string
	Zone     *string
	Proxy    *string
	File     *string
	Type     ProviderType
	Supplier supplier.Supplier
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
