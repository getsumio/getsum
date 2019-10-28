package servers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	. "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/logger"
	"github.com/getsumio/getsum/internal/status"
	. "github.com/getsumio/getsum/internal/supplier"
	"github.com/getsumio/getsum/internal/validation"
)

type OnPremiseServer struct {
	address     string
	port        int
	storagePath string
	supplier    Supplier
	mux         sync.Mutex
}

var factory ISupplierFactory

func (s *OnPremiseServer) Start() {
	logger.Level = logger.LevelInfo
	http.HandleFunc("/", s.handle)
	listenAddress := fmt.Sprintf("%s:%d", s.address, s.port)
	http.ListenAndServe(listenAddress, nil)
}

func (s *OnPremiseServer) handle(w http.ResponseWriter, r *http.Request) {
	s.mux.Lock()
	defer s.mux.Unlock()

	logger.Info("Request received on %s with method %s", r.Method, r.URL.Path)
	switch r.Method {
	case "GET":
		if s.supplier == nil {
			handleError("There is no running process", w)
			return
		}
		status, _ := json.Marshal(s.supplier.Status())
		w.Write(status)
	case "POST":
		if s.supplier != nil {
			handleError("Server already running a process", w)
			return
		}
		jsonDecoder := json.NewDecoder(r.Body)
		config := &Config{}
		err := jsonDecoder.Decode(config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = validation.ValidateConfig(config)
		if err != nil {
			handleError(err.Error(), w)
			return

		}

		if factory == nil {
			factory = new(SupplierFactory)
		}

		var algorithm = ValueOf(&config.Algorithm[0])
		supplier := factory.GetSupplierByAlgo(config, &algorithm)
		s.supplier = supplier
		go s.supplier.Run()
	case "DELETE":
		if s.supplier == nil {
			handleError("There is no running process", w)
			return
		}
		s.supplier.Terminate()
		w.WriteHeader(http.StatusOK)
		s.supplier = nil
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		logger.Error("Can not handle request method rejecting request")
	}

}

func handleError(message string, w http.ResponseWriter) {
	status := &status.Status{Type: status.ERROR, Value: message}
	jsonStatus, _ := json.Marshal(status)
	w.Write(jsonStatus)
	return
}
