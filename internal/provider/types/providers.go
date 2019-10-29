package providers

import (
	"sync"
	"time"

	"github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/status"
)

type Providers struct {
	List     []Provider
	config   *config.Config
	quit     chan bool
	wait     chan bool
	channels []<-chan *status.Status
	statuses []*status.Status
	length   int
}

func (p *Providers) init() {
	p.length = p.Size()
	p.quit, p.wait = make(chan bool, p.length), make(chan bool)
	p.statuses = make([]*status.Status, p.length)
	p.channels = make([]<-chan *status.Status, p.length)
}

func (p *Providers) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	p.init()

}

func (p *Providers) Terminate() {
	for i := 0; i < p.length; i++ {
		p.quit <- true
	}
	close(p.quit)
	time.Sleep(200 * time.Millisecond)

}

func (p *Providers) HasError() bool {
	return false
}

func (p *Providers) IsFinish() bool {
	return false
}

func (p *Providers) Size() int {
	if p.List == nil {
		return 0
	}
	return len(p.List)
}
