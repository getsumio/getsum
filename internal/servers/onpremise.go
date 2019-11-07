package servers

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/getsumio/getsum/internal/config"
	. "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	"github.com/getsumio/getsum/internal/status"
	. "github.com/getsumio/getsum/internal/supplier"
	"github.com/getsumio/getsum/internal/validation"
	"github.com/google/uuid"
)

//server instance to run in server mode
type OnPremiseServer struct {
	StoragePath string
	mux         *sync.RWMutex
	suppliers   map[string]Supplier
}

var factory ISupplierFactory

const uuidPattern = "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"

var regex *regexp.Regexp = regexp.MustCompile(uuidPattern)

const default_capacity = 250

//start server in given config listen address and port or tls details
func (s *OnPremiseServer) Start(config *config.Config) error {
	logger.Level = logger.LevelInfo
	factory = new(SupplierFactory)
	s.suppliers = make(map[string]Supplier)
	s.mux = &sync.RWMutex{}
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
func handleGet(s *OnPremiseServer, w http.ResponseWriter, r *http.Request, id string) {
	//check if any runner
	s.mux.RLock()
	defer s.mux.RUnlock()
	logger.Info("Checking if there is a process with id : %s", id)
	supplier, ok := s.suppliers[id]
	if !ok {
		handleError("There is no running process", w)
		return
	}
	//collect status and return
	stat := supplier.Status()
	status, err := json.Marshal(stat)
	if err != nil {
		handleError("System can not parse given status %s", w, err.Error())
		return
	}
	logger.Info("Returning status %v", *stat)
	w.Write(status)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "application/json")

}

//post executed to Run a new calculation
func handlePost(s *OnPremiseServer, w http.ResponseWriter, r *http.Request) {
	s.mux.Lock()
	defer s.mux.Unlock()
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
	supplier, err := factory.GetSupplierByAlgo(config, &algorithm)
	if err != nil {
		handleError("Can not create algorithm runner instance: "+err.Error(), w)
		return
	}
	processId := uuid.New().String()
	stat := supplier.Status()
	stat.Type = status.STARTED
	stat.Value = processId
	jsonStat, err := json.Marshal(*stat)
	if err != nil {
		handleError("Can not parse status"+err.Error(), w)
		return
	}
	w.Write(jsonStat)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,POST,GET,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type,accept,origin")

	supplier.Data().StartTime = time.Now()
	go supplier.Run()
	s.suppliers[processId] = supplier
	logger.Info("Process started")
	s.ensureCapacity()
}

func handleOptions(s *OnPremiseServer, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,POST,GET,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type,accept,origin")
	w.Header().Set("Allow", "OPTIONS,GET,POST,DELETE")

}

//terminates running calculation
func handleDelete(s *OnPremiseServer, w http.ResponseWriter, r *http.Request, id string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	supplier, ok := s.suppliers[id]
	if !ok {
		handleError("There is no running process", w)
		return
	}

	supplier.Terminate()
	supplier.Delete()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,POST,GET,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type,accept,origin")

	delete(s.suppliers, id)
	logger.Info("Process terminated")
}

//delegates GET POST DELETE main server listener
func (s *OnPremiseServer) handle(w http.ResponseWriter, r *http.Request) {
	logger.LogRequest(r)
	switch r.Method {
	case "GET":
		requestId := path.Base(html.EscapeString(r.URL.Path))
		if !regex.MatchString(requestId) {
			logger.Error("Request is not a valid request id! %s", requestId)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		handleGet(s, w, r, requestId)
	case "POST":
		handlePost(s, w, r)
	case "OPTIONS":
		handleOptions(s, w, r)
	case "DELETE":
		requestId := path.Base(html.EscapeString(r.URL.Path))
		if !regex.MatchString(requestId) {
			logger.Error("Request is not a valid request id! %s", requestId)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		handleDelete(s, w, r, requestId)
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

func (s *OnPremiseServer) ensureCapacity() {
	if len(s.suppliers) >= default_capacity {
		now := time.Now()
		for k := range s.suppliers {
			supplier := s.suppliers[k]
			if int(now.Sub(supplier.Data().StartTime).Seconds()) > (supplier.Data().TimeOut * 2) {
				delete(s.suppliers, k)
			}
		}
	}
}
