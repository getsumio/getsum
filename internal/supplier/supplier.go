package supplier

import (
	"time"

	"github.com/getsumio/getsum/internal/file"
	"github.com/getsumio/getsum/internal/status"
)

//Supplier ares main runner for calculation
type Supplier interface {
	Run()
	Status() *status.Status
	Terminate() error
	Supports() []Algorithm
	Delete()
	Data() *BaseSupplier
}

//embedded struct for suppliers
type BaseSupplier struct {
	Algorithm    Algorithm
	status       *status.Status
	File         *file.File
	TimeOut      int
	Key          string
	StartTime    time.Time
	IsConcurrent bool
}
