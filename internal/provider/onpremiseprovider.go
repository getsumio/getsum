package providers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider/types"
	"github.com/getsumio/getsum/internal/status"
)

type RemoteProvider struct {
	BaseProvider
	client      *http.Client
	config      *config.Config
	address     string
	ErrorStatus *status.Status
}

func (l *RemoteProvider) Run() {
	if l.BaseProvider.Wait {
		logger.Info("Process %s on hold", l.Type.Name())
		l.WG.Wait()
	}
	logger.Debug("Running remote provider %s", l.Name)
	remoteRun(l)
}

func getErrorStatus(err error) *status.Status {
	stat := &status.Status{}
	stat.Type = status.ERROR
	stat.Value = err.Error()
	return stat
}

func remoteRun(l *RemoteProvider) {

	body, err := json.Marshal(*l.config)
	if err != nil {
		l.ErrorStatus = getErrorStatus(err)
		return
	}

	resp, err := l.client.Post(l.address, "application/json", bytes.NewBuffer(body))
	defer closeResponse(resp)
	if err != nil {
		l.ErrorStatus = getErrorStatus(err)
		return

	}

}

func closeResponse(response *http.Response) {
	if response != nil && response.Body != nil {
		response.Body.Close()
	}
}

func remoteStatus(l *RemoteProvider) *status.Status {
	resp, err := l.client.Get(l.address)
	if err != nil {
		return getErrorStatus(err)
	}

	defer closeResponse(resp)
	decoder := json.NewDecoder(resp.Body)
	status := &status.Status{}
	err = decoder.Decode(status)
	if err != nil {
		return getErrorStatus(err)
	}

	return status
}

func remoteTerminate(l *RemoteProvider) error {
	req, err := http.NewRequest("DELETE", l.address, nil)
	if err != nil {
		return err
	}

	resp, err := l.client.Do(req)
	defer closeResponse(resp)
	if err != nil {
		return err
	}

	return nil
}

func (l *RemoteProvider) Data() *BaseProvider {
	return &l.BaseProvider
}

func (l *RemoteProvider) Wait() {
	logger.Info("Provider %s suspended", l.Name)
	l.BaseProvider.Wait = true
	l.WG.Add(1)
	stat := l.Supplier.Status()
	stat.Type = status.SUSPENDED
}

func (l *RemoteProvider) Resume() {
	logger.Info("Resuming %s", l.Name)
	l.WG.Done()
	stat := l.Supplier.Status()
	stat.Type = status.RESUMING
}

func (l *RemoteProvider) Terminate() error {
	logger.Debug("Quit triggered %s", l.Name)
	return remoteTerminate(l)

}

func (l *RemoteProvider) Status() *status.Status {
	var stat *status.Status
	if l.ErrorStatus != nil {
		stat = l.ErrorStatus
	} else {
		stat = remoteStatus(l)
	}

	logger.Trace("%s - Remote status received: %v", l.Name, *stat)
	return stat
}
