package providers

import (
	"time"

	. "github.com/getsumio/getsum/internal/file"
	. "github.com/getsumio/getsum/internal/provider/types"
)

type LocalProvider struct {
	BaseProvider
}

func (l *LocalProvider) Run(quit <-chan bool, wait <-chan bool) <-chan *Status {
	defer complete(l)
	statusChannel := make(chan *Status)
	go l.Supplier.Run()
	go func() {
		for {
			select {
			case <-wait:
			case <-quit:
				l.Supplier.Terminate()
				complete(l)
				return
			default:
				stat := l.Supplier.Status()
				statusChannel <- stat
				time.Sleep(50 * time.Millisecond)

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
}
