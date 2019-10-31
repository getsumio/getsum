package providers

import (
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider/types"
	"github.com/getsumio/getsum/internal/status"
)

//local provider runs on host pc
type LocalProvider struct {
	BaseProvider
}

//starts calculating checksum
//according to user param
//related supplier (os,openssl or os apps) will be called
func (l *LocalProvider) Run() {
	//wait if suspend called
	if l.BaseProvider.Wait {
		logger.Info("Provider Process on hold! %s", l.Type.Name())
		l.WG.Wait()
	}

	logger.Debug("Running local provider %s", l.Name)
	l.Supplier.Run()
}

//a shortcut access to embedded struct
//in case of interface used as reference Value
func (l *LocalProvider) Data() *BaseProvider {
	return &l.BaseProvider
}

//This should be called
//before Run method
//marks provider as wait
//so runner will refer to WaitGroup first
//instead of running
func (l *LocalProvider) Wait() {
	logger.Info("Provider %s suspended", l.Name)
	l.BaseProvider.Wait = true
	l.WG.Add(1)
	stat := l.Supplier.Status()
	stat.Type = status.SUSPENDED
}

//if process suspended
//notifies pending routine to continue
func (l *LocalProvider) Resume() {
	logger.Info("Resuming %s", l.Name)
	l.BaseProvider.Wait = false
	l.WG.Done()
	stat := l.Supplier.Status()
	stat.Type = status.RESUMING
}

//calls terminate on running calculation
//how terminate handled may vary according to lib parameter
//if its openssl or unix command it will call process kill
//but if go library it will just update status
func (l *LocalProvider) Terminate() error {
	logger.Debug("Quit triggered %s", l.Name)
	return l.Supplier.Terminate()
}

//collect status from running process and returns
func (l *LocalProvider) Status() *status.Status {
	stat := l.Supplier.Status()
	logger.Trace("%s Returning status %v address %v supplier: %v ", l.Name, *stat, stat, l.Supplier)
	return stat
}

//remove file from local pc
//will be called if server mode or validation fails
func (l *LocalProvider) DeleteFile() {
	logger.Debug("Removing file %s", l.Name)
	l.Supplier.Delete()
}
