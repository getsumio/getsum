package supplier

import (
	"github.com/getsumio/getsum/internal/file"
)

type Supplier interface {
	Run()
	Status() *Status
	Terminate()
}

type BaseSupplier struct {
	Algorithm Algorithm
	status    *Status
	File      *file.File
	TimeOut   int
}

type Status struct {
	Status   string
	Value    string
	Checksum string
}
