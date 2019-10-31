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

//reaches given server
//using http client
//and collect status
//or run/terminates
type RemoteProvider struct {
	BaseProvider
	client      *http.Client
	config      *config.Config
	address     string
	status      *status.Status
	hasRunError bool
}

//notifies server to run
//waits if process suspended
func (l *RemoteProvider) Run() {
	if l.BaseProvider.Wait {
		logger.Info("Process %s on hold", l.Type.Name())
		l.WG.Wait()
	}
	logger.Debug("Running remote provider %s", l.Name)
	remoteRun(l)
}

//utility to create a status struct with given value
func setErrorStatus(err error, r *RemoteProvider) {
	r.status.Type = status.ERROR
	r.status.Value = err.Error()
}

//send request to server to start running
func remoteRun(l *RemoteProvider) {

	l.status = &status.Status{}

	//send config to server, this only POST request so no need other param
	resp, err := l.client.Post(l.address, "application/json", bytes.NewBuffer(config.ConfigJson))
	defer closeResponse(resp)
	if err != nil {
		//set error as provider status
		//status() method will handle
		setErrorStatus(err, l)
		l.hasRunError = true
		return

	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(l.status)
	if err != nil {
		setErrorStatus(err, l)
		l.hasRunError = true
	}
	l.hasRunError = l.status.Type == status.ERROR

}

//utility to close response if present
func closeResponse(response *http.Response) {
	if response != nil && response.Body != nil {
		response.Body.Close()
	}
}

//fetches server using GET and collects its status
func remoteStatus(l *RemoteProvider) *status.Status {
	//reach the server
	resp, err := l.client.Get(l.address)
	if err != nil {
		setErrorStatus(err, l)
		return l.status
	}

	defer closeResponse(resp)
	//parse response
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(l.status)
	if err != nil {
		setErrorStatus(err, l)
		return l.status
	}

	return l.status
}

//trigger termination on remote server using http DELETE
func remoteTerminate(l *RemoteProvider) error {
	//let the server know process terminated
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

//shorthand to embedded struct in case of interface used
func (l *RemoteProvider) Data() *BaseProvider {
	return &l.BaseProvider
}

//suspend this runner
//SHOULD BE CALLED BEFORE Run method
func (l *RemoteProvider) Wait() {
	logger.Info("Provider %s suspended", l.Name)
	l.BaseProvider.Wait = true
	l.WG.Add(1)
	stat := l.Supplier.Status()
	stat.Type = status.SUSPENDED
}

//resume this provider
func (l *RemoteProvider) Resume() {
	logger.Info("Resuming %s", l.Name)
	l.WG.Done()
	stat := l.Supplier.Status()
	stat.Type = status.RESUMING
}

//triggers terminate on remote server
func (l *RemoteProvider) Terminate() error {
	logger.Debug("Quit triggered %s", l.Name)
	return remoteTerminate(l)

}

//no matter what remote server always deletes file
//this is interface impl
func (l *RemoteProvider) DeleteFile() {
	//Do nothing
}

//collect status from remote server
func (l *RemoteProvider) Status() *status.Status {
	var stat *status.Status
	//check if this provided already faced an error
	//if so dont bother raching to server
	if l.hasRunError {
		stat = l.status
	} else {
		stat = remoteStatus(l)
	}

	logger.Trace("%s - Remote status received: %v", l.Name, *stat)
	return stat
}
