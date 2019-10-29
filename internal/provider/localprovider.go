package providers

import (
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider/types"
	"github.com/getsumio/getsum/internal/status"
)

type LocalProvider struct {
	BaseProvider
}

func (l *LocalProvider) Run() {
	if l.BaseProvider.Wait {
		logger.Info("Process on hold! %s", l.Type.Name())
		l.WG.Wait()
	}
	logger.Debug("Running local provider %s", l.Name)
	l.Supplier.Run(false)
}

func (l *LocalProvider) Data() *BaseProvider {
	return &l.BaseProvider
}

func (l *LocalProvider) Wait() {
	logger.Info("Provider %s suspended", l.Name)
	l.BaseProvider.Wait = true
	l.WG.Add(1)
	stat := l.Supplier.Status()
	stat.Type = status.SUSPENDED
}

func (l *LocalProvider) Resume() {
	logger.Info("Resuming %s", l.Name)
	l.WG.Done()
	stat := l.Supplier.Status()
	stat.Type = status.RESUMING
}

func (l *LocalProvider) Terminate() error {
	logger.Debug("Quit triggered %s", l.Name)
	return l.Supplier.Terminate()
}

func (l *LocalProvider) Status() *status.Status {
	return l.Supplier.Status()
}
