package providers

import (
	"time"

	. "github.com/getsumio/getsum/internal/algorithm/supplier"
)

type LocalProvider struct {
	BaseProvider
}

func (l *LocalProvider) Run(quit <-chan bool, wait <-chan bool) <-chan *Status {
	l.Status = "STARTED"
	defer complete(l)
	statusChannel := make(chan *Status)
	l.Supplier.Run()
	go func() {
		for {
			select {
			case <-wait:
			case <-quit:
				l.Supplier.Terminate()
				complete(l)
				return
			default:
				statusChannel <- l.Supplier.Status()
				time.Sleep(500 * time.Millisecond)

			}
		}
	}()
	return statusChannel
}

func (l *LocalProvider) Data() *BaseProvider {
	return &l.BaseProvider
}

func (l *LocalProvider) Close() {

}

func complete(l *LocalProvider) {
	l.Status = "COMPLETED"
}
