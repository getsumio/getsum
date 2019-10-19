package supplier

import (
	"fmt"
)

type UnixSupplier struct {
	BaseSupplier
}

func (s UnixSupplier) Run() (string, error) {
	return fmt.Sprintf("Selected %d", s.Algorithm), nil
}
