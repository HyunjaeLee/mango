package parse

import (
	"sync"
	"fmt"
)

type Pipe struct {
	request chan func()
	done chan struct{}
	wg *sync.WaitGroup
}

func NewPipe() *Pipe {
	return &Pipe{
		request: make(chan func()),
		done: make(chan struct{}),
		wg: new(sync.WaitGroup),
	}
}

func (p *Pipe) Worker() {
	for req := range p.request {
		req()
	}
}

func (p *Pipe) Start() {
	const numWorkers = 10
	p.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			p.Worker()
			fmt.Println("work!")
			p.wg.Done()
		}()
	}
}

func (p *Pipe) Execute(req func()) {
	p.request <- req
}

func (p *Pipe) Wait() {
	p.wg.Wait()
}

func (p *Pipe) Close() {
	close(p.request)
	close(p.done)
}