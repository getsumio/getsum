package providers

import (
	. "github.com/getsumio/getsum/internal/algorithm/supplier"
)

type Provider interface {
	Close()
	Data() *BaseProvider
	Run() (string, error)
}

type BaseProvider struct {
	Status   string
	Value    string
	Name     string
	Address  string
	Zone     string
	Proxy    string
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
