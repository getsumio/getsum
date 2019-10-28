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

	body, _ := json.Marshal(*l.config)
	resp, _ := l.client.Post(l.address, "application/json", bytes.NewBuffer(body))
	defer resp.Body.Close()
}

func remoteStatus(l *RemoteProvider) *status.Status {
	resp, _ := l.client.Get(l.address)
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	status := &status.Status{}
	decoder.Decode(status)
	return status
}

func remoteTerminate(l *RemoteProvider) {
	req, err := http.NewRequest("DELETE", l.address, nil)
	if err != nil {
		panic(err)
	}

	resp, _ := l.client.Do(req)
	defer resp.Body.Close()
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
