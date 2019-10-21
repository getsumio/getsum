package supplier

import (
	"fmt"
	"time"
)

type UnixSupplier struct {
	BaseSupplier
}

func (s *UnixSupplier) Run() error {
	t := time.After(20 * time.Second)
	s.status = &Status{"STARTED", ""}
	var i int
	go func() {
		for {
			select {
			case <-t:
				s.status = &Status{"COMPLETED", fmt.Sprintf("%d%%", i)}
				return
			default:
				i += 5
				val := fmt.Sprintf("%d%%", i)
				s.status = &Status{"RUNNING", val}
				time.Sleep(time.Second)
			}
		}
	}()
	return nil
}

func (s *UnixSupplier) Status() *Status {
	return s.status
}

func (s *UnixSupplier) Terminate() error {
	s.status = &Status{"Terminated", ""}
	return nil
}
