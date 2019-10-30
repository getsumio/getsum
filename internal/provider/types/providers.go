package providers

import (
	"sync"
	"time"

	"github.com/getsumio/getsum/internal/status"
)

//TODO use single loop instead for each method
type Providers struct {
	Remotes       []*Provider
	Locals        []*Provider
	All           []*Provider
	Statuses      []*status.Status
	Length        int
	HasRemote     bool
	HasLocal      bool
	HasValidation bool
	mux           sync.Mutex
}

func (providers *Providers) RunRemotes() {
	for _, provider := range providers.Remotes {
		go (*provider).Run()
	}
	time.Sleep(200 * time.Millisecond)
}

func (providers *Providers) RunLocals() {
	for _, provider := range providers.Locals {
		go (*provider).Run()
	}
}

func (providers *Providers) Run() {
	providers.RunRemotes()
	providers.RunLocals()
}

func (providers *Providers) SuspendLocals() {
	for _, provider := range providers.Locals {
		(*provider).Wait()
	}
}

func (providers *Providers) Delete() {
	for _, p := range providers.All {
		(*p).DeleteFile()
	}
}

func (providers *Providers) ResumeLocals() {
	for _, provider := range providers.Locals {
		(*provider).Resume()
	}
}

func (providers *Providers) Terminate(force bool) {
	for i, provider := range providers.All {
		if force || providers.Statuses[i].Type >= status.COMPLETED {
			(*provider).Terminate()
		}
	}
}

func (providers *Providers) Status() []*status.Status {
	providers.mux.Lock()
	defer providers.mux.Unlock()
	for i, provider := range providers.All {
		if providers.Statuses[i] == nil || providers.Statuses[i].Type < status.COMPLETED {
			providers.Statuses[i] = (*provider).Status()
		}
	}
	return providers.Statuses
}

func (providers *Providers) HasError() bool {
	for _, stat := range providers.Statuses {
		if stat.Type > status.COMPLETED {
			return true
		}
	}
	return false
}

func (providers *Providers) IsRunning() bool {
	for _, stat := range providers.Status() {
		if stat.Type < status.COMPLETED {
			return true
		}
	}
	return false
}

func (providers *Providers) HasMismatch(checksum string) bool {

	if checksum == "" {
		return false
	}
	providers.mux.Lock()
	defer providers.mux.Unlock()
	var mismatch bool = false
	for i, stat := range providers.Statuses {
		if stat.Type == status.COMPLETED {
			if stat.Checksum != checksum {
				providers.Statuses[i].Type = status.MISMATCH
				mismatch = true
			} else {
				providers.Statuses[i].Type = status.VALIDATED
			}
		}
	}
	return mismatch

}
