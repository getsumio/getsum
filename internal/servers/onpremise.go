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

//server instance to run in server mode
type OnPremiseServer struct {
	StoragePath string
	Supplier    Supplier
	mux         sync.Mutex
}

var factory ISupplierFactory

//start server in given config listen address and port or tls details
//TODO add interface support
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

//get executed to reach status
//so collect status if there is any runner
func handleGet(s *OnPremiseServer, w http.ResponseWriter, r *http.Request) {
	//check if any runner
	if s.Supplier == nil {
		handleError("There is no running process", w)
		return
	}
	//collect status and return
	stat := s.Supplier.Status()
	status, err := json.Marshal(stat)
	if err != nil {
		handleError("System can not parse given status %s", w, err.Error())
		return
	}
	logger.Info("Returning status %v", *stat)
	w.Write(status)
}

//post executed to Run a new calculation
func handlePost(s *OnPremiseServer, w http.ResponseWriter, r *http.Request) {
	//check if any runner
	if s.Supplier != nil {
		stat := s.Supplier.Status()
		if stat.Type <= status.RUNNING {
			handleError("Server already running another process", w)
			return
		}
	}
	//read the config from request
	jsonDecoder := json.NewDecoder(r.Body)
	config := &Config{}
	err := jsonDecoder.Decode(config)
	if err != nil {
		handleError("Can not read given config %s", w, err.Error())
		return
	}
	//validate config
	err = validation.ValidateConfig(config, true)
	if err != nil {
		handleError(err.Error(), w)
		return

	}

	//get supplier instance, only single algo supported on server mode
	var algorithm = ValueOf(&config.Algorithm[0])
	config.Dir = &s.StoragePath
	s.Supplier, err = factory.GetSupplierByAlgo(config, &algorithm)
	//something went  wrong, TODO add error handler
	if err != nil {
		handleError("Can not create algorithm runner instance: "+err.Error(), w)
		return
	}
	go s.Supplier.Run()
	handleGet(s, w, r)
	logger.Info("Process started")
}

//terminates running calculation
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

//delegates GET POST DELETE main server listener
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

//utility to write given error to response
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
