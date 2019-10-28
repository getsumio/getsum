package providers

import (
	"net/http"
	"time"

	"github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider/types"
	"github.com/getsumio/getsum/internal/status"
)

type RemoteProvider struct {
	BaseProvider
	client http.Client
	config *config.Config
}

func (l *RemoteProvider) Run(quit <-chan bool, wait <-chan bool) <-chan *status.Status {
	logger.Debug("Running remote provider %s", l.Name)
	defer l.Close()
	statusChannel := make(chan *status.Status)
	logger.Trace("Triggering supplier %s", l.Name)
	go remoteRun(l)
	go func() {
		for {
			select {
			case <-wait:
			case <-quit:
				logger.Trace("Quit triggered %s", l.Name)
				remoteTerminate(l)
				return
			default:
				stat := remoteStatus(l)
				logger.Trace("Status received", (*stat).Type.Name(), (*stat).Value, l.Name)
				statusChannel <- stat
				time.Sleep(50 * time.Millisecond)

			}
		}
	}()
	return statusChannel
}

func remoteRun(l *RemoteProvider) {

}

func remoteStatus(l *RemoteProvider) *status.Status {
	return nil
}

func remoteTerminate(l *RemoteProvider) {

}

func (l *RemoteProvider) Data() *BaseProvider {
	return &l.BaseProvider
}

func (l *RemoteProvider) Close() {
	l.Supplier.Terminate()
}

func (l *RemoteProvider) Region() string {
	return l.Zone
}
