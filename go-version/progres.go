package main

import (
	"sync"

	"github.com/pterm/pterm"
)

type Progress struct {
	bar  *pterm.ProgressbarPrinter
	lock sync.Mutex
}

func NewProgress(total int) Progress {
	bar, err := pterm.DefaultProgressbar.WithTotal(int(total)).WithTitle("First test").Start()
	if err != nil {
		panic(err)
	}

	return Progress{bar: bar}
}

func (p *Progress) Update(value int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.bar.Add(value)
}

func (p *Progress) Increment() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.bar.Increment()
}
