package providers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	. "github.com/getsumio/getsum/internal/provider/types"
	"github.com/getsumio/getsum/internal/status"
)

type RemoteProvider struct {
	BaseProvider
	client  *http.Client
	config  *config.Config
	address string
}

func (l *RemoteProvider) Run(quit <-chan bool, wait <-chan bool) <-chan *status.Status {
	logger.Debug("Running remote provider %s", l.Name)
	statusChannel := make(chan *status.Status)
	logger.Trace("Triggering supplier %s", l.Name)
	go remoteRun(l, statusChannel)
	go func() {
		time.Sleep(time.Second)
		for {
			select {
			case <-wait:
			case <-quit:
				logger.Trace("Quit triggered %s", l.Name)
				err := remoteTerminate(l)
				if err != nil {
					statusChannel <- getErrorStatus(err)
				}
				return
			default:
				stat := remoteStatus(l)
				logger.Trace("Status received", (*stat).Type.Name(), (*stat).Value, l.Name)
				statusChannel <- stat
				time.Sleep(150 * time.Millisecond)

			}
		}
	}()
	return statusChannel
}

func getErrorStatus(err error) *status.Status {
	stat := &status.Status{}
	stat.Type = status.ERROR
	stat.Value = err.Error()
	return stat
}

func remoteRun(l *RemoteProvider, statusChannel chan *status.Status) {

	body, err := json.Marshal(*l.config)
	if err != nil {
		statusChannel <- getErrorStatus(err)
		return
	}

	resp, err := l.client.Post(l.address, "application/json", bytes.NewBuffer(body))
	defer resp.Body.Close()

	if err != nil {
		statusChannel <- getErrorStatus(err)
		return

	}

}

func remoteStatus(l *RemoteProvider) *status.Status {
	resp, err := l.client.Get(l.address)
	if err != nil {
		return getErrorStatus(err)
	}

	defer resp.Body.Close()
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
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (l *RemoteProvider) Data() *BaseProvider {
	return &l.BaseProvider
}

func (l *RemoteProvider) Close() {
	remoteTerminate(l)
}

func (l *RemoteProvider) Region() string {
	return l.Zone
}
