package providers

import (
	"sync"

	"github.com/getsumio/getsum/internal/status"
	"github.com/getsumio/getsum/internal/supplier"
)

type Provider interface {
	Data() *BaseProvider
	Run()
	Wait()
	Resume()
	Status() *status.Status
	Terminate() error
	DeleteFile()
}

type BaseProvider struct {
	Name     string
	Type     ProviderType
	Supplier supplier.Supplier
	WG       *sync.WaitGroup
	Wait     bool
}

type ProviderType int

const (
	Aws ProviderType = iota
	Google
	Oracle
	Azure
	IBM
	Local
	OnPremise
)

var typStr = []string{
	"aws",
	"google",
	"oracle",
	"azure",
	"ibm",
	"local",
	"onpremise",
}

func (p ProviderType) Name() string {
	return typStr[p]
}
