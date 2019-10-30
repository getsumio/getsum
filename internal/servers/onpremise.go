package servers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/getsumio/getsum/internal/config"
	. "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	"github.com/getsumio/getsum/internal/status"
	. "github.com/getsumio/getsum/internal/supplier"
	"github.com/getsumio/getsum/internal/validation"
)

type OnPremiseServer struct {
	storagePath string
	Supplier    Supplier
	mux         sync.Mutex
}

var factory ISupplierFactory

func (s *OnPremiseServer) Start(config *config.Config) error {
	logger.Level = logger.LevelInfo
	factory = new(SupplierFactory)
	http.HandleFunc("/", s.handle)
	listenAddress := fmt.Sprintf("%s:%d", *config.Listen, *config.Port)
	var err error
	if *config.TLSKey != "" {
		err = http.ListenAndServeTLS(listenAddress, *config.TLSCert, *config.TLSKey, nil)
	} else {
		err = http.ListenAndServe(listenAddress, nil)
	}
	if err != nil {
		return err
	}

	return nil
}

func handleGet(s *OnPremiseServer, w http.ResponseWriter, r *http.Request) {
	if s.Supplier == nil {
		handleError("There is no running process", w)
		return
	}
	stat := s.Supplier.Status()
	status, err := json.Marshal(stat)
	if err != nil {
		handleError("System can not parse given status %s", w, err.Error())
		return
	}
	logger.Info("Returning status %v", *stat)
	w.Write(status)
}

func handlePost(s *OnPremiseServer, w http.ResponseWriter, r *http.Request) {
	if s.Supplier != nil {
		stat := s.Supplier.Status()
		if stat.Type <= status.RUNNING {
			handleError("Server already running another process", w)
			return
		}
	}
	jsonDecoder := json.NewDecoder(r.Body)
	config := &Config{}
	err := jsonDecoder.Decode(config)
	if err != nil {
		handleError("Can not read given config %s", w, err.Error())
		return
	}
	err = validation.ValidateConfig(config, true)
	if err != nil {
		handleError(err.Error(), w)
		return

	}

	var algorithm = ValueOf(&config.Algorithm[0])
	s.Supplier = factory.GetSupplierByAlgo(config, &algorithm)
	if s.Supplier == nil {
		handleError("Can not create algorithm runner instance", w)
	}
	go s.Supplier.Run()
	logger.Info("Process started")
}

func handleDelete(s *OnPremiseServer, w http.ResponseWriter, r *http.Request) {
	if s.Supplier == nil {
		handleError("There is no running process", w)
		return
	}
	s.Supplier.Terminate()
	s.Supplier.Delete()
	w.WriteHeader(http.StatusOK)
	s.Supplier = nil
	logger.Info("Process terminated")
}

func (s *OnPremiseServer) handle(w http.ResponseWriter, r *http.Request) {
	s.mux.Lock()
	defer s.mux.Unlock()

	logger.LogRequest(r)
	switch r.Method {
	case "GET":
		handleGet(s, w, r)
	case "POST":
		handlePost(s, w, r)
	case "DELETE":
		handleDelete(s, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		logger.Error("Can not handle request method rejecting request")
	}

}

func handleError(message string, w http.ResponseWriter, params ...interface{}) {
	if params != nil {
		message = fmt.Sprintf(message, params...)
	}
	logger.Error("An unexpected error occured %s", message)

	status := &status.Status{Type: status.ERROR, Value: message}
	jsonStatus, _ := json.Marshal(status)
	w.Write(jsonStatus)
	return
}
