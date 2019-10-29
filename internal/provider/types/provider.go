package providers

import (
	"github.com/getsumio/getsum/internal/status"
	"github.com/getsumio/getsum/internal/supplier"
)

type Provider interface {
	Data() *BaseProvider
	Run(quit <-chan bool, wait <-chan bool) <-chan *status.Status
}

type BaseProvider struct {
	Name     string
	Proxy    *string
	Type     ProviderType
	Supplier supplier.Supplier
}

type ProviderType int

const (
	Aws    ProviderType = iota
	Google ProviderType = iota
	Oracle ProviderType = iota
	Azure  ProviderType = iota
	IBM    ProviderType = iota
	Local  ProviderType = iota
)

var typStr = []string{
	"AWS",
	"GOOGLE",
	"ORACLE",
	"AZURE",
	"IBM",
	"LOCAL",
}

func (p ProviderType) Name() string {
	return typStr[p]
}
