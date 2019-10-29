package providers

import (
	"time"

	"github.com/getsumio/getsum/internal/status"
)

type Providers struct {
	Remotes   []*Provider
	Locales   []*Provider
	All       []*Provider
	Statuses  []*status.Status
	Length    int
	HasRemote bool
	HasLocal  bool
}

func (providers *Providers) RunRemotes() {
	for _, provider := range providers.Remotes {
		go (*provider).Run()
	}
	time.Sleep(200 * time.Millisecond)
}

func (providers *Providers) RunLocales() {
	for _, provider := range providers.Locales {
		go (*provider).Run()
	}
}

func (providers *Providers) Run() {
	providers.RunRemotes()
	providers.RunLocales()
}

func (providers *Providers) SuspendLocales() {
	for _, provider := range providers.Locales {
		(*provider).Wait()
	}
}

func (providers *Providers) ResumeLocales() {
	for _, provider := range providers.Locales {
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
	var mismatch bool = false
	for _, stat := range providers.Statuses {
		if stat.Type == status.COMPLETED && stat.Checksum != checksum {
			stat.Type = status.MISMATCH
			mismatch = true
		}
	}
	return mismatch

}
