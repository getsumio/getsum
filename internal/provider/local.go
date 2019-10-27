package providers

import (
	"time"

	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider/types"
	"github.com/getsumio/getsum/internal/status"
)

type LocalProvider struct {
	BaseProvider
}

func (l *LocalProvider) Run(quit <-chan bool, wait <-chan bool) <-chan *status.Status {
	logger.Debug("Running local provider %s", l.Name)
	defer l.Close()
	statusChannel := make(chan *status.Status)
	logger.Trace("Triggering supplier %s", l.Name)
	go l.Supplier.Run()
	go func() {
		for {
			select {
			case <-wait:
			case <-quit:
				logger.Trace("Quit triggered %s", l.Name)
				l.Supplier.Terminate()
				return
			default:
				stat := l.Supplier.Status()
				logger.Trace("Status received", (*stat).Type.Name(), (*stat).Value, l.Name)
				statusChannel <- stat
				time.Sleep(50 * time.Millisecond)

			}
		}
	}()
	return statusChannel
}

func (l *LocalProvider) Data() *BaseProvider {
	return &l.BaseProvider
}

func (l *LocalProvider) Close() {
	l.Supplier.Terminate()
}

func (l *LocalProvider) Region() string {
	return *l.Zone
}
