package servers

import (
	"github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/status"
)

type GetsumRequest struct {
	RequestType string        `json:"request_type"`
	Config      config.Config `json:"config"`
}

type Server interface {
	Handle(request GetsumRequest) (status.Status, error)
}
