package supplier

import (
	"github.com/getsumio/getsum/internal/file"
	"github.com/getsumio/getsum/internal/status"
)

type Supplier interface {
	Run()
	Status() *status.Status
	Terminate() error
	Supports() []Algorithm
	Delete()
}

type BaseSupplier struct {
	Algorithm Algorithm
	status    *status.Status
	File      *file.File
	TimeOut   int
	Key       string
}
