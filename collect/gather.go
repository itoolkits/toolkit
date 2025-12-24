// gather use to batch collect

package collect

import (
	"time"
)

type Gather[T any] struct {
	timeout time.Duration
	size    int
	call    func([]T, *Gather[T])

	ch chan T

	batch []T

	done chan struct{}

	tck *time.Ticker
}

// NewGather create gather struct, use one goroutine
func NewGather[T any](timeout time.Duration, size int, call func([]T, *Gather[T])) *Gather[T] {
	g := &Gather[T]{
		timeout: timeout,
		size:    size,
		ch:      make(chan T, size*2),
		batch:   make([]T, 0, size),
		call:    call,
		tck:     time.NewTicker(timeout),
		done:    make(chan struct{}),
	}

	// goroutine
	go g.init()

	return g
}

// init init struct
func (g *Gather[T]) init() {
	defer func() {
		g.tck.Stop()
	}()

	if g.size <= 1 {
		for v := range g.ch {
			g.call([]T{v}, g)
		}
		return
	}

	for {
		select {
		case ele, ok := <-g.ch:
			if !ok {
				g.callback()
				close(g.done)
				g.tck.Stop()
				return
			}
			g.batch = append(g.batch, ele)
			if len(g.batch) >= g.size {
				g.callback()
				g.resetTicker()
			}
		case <-g.tck.C:
			g.callback()
		}
	}
}

// resetTicker reset ticker
func (g *Gather[T]) resetTicker() {
	g.tck.Reset(g.timeout)
}

// Close - close gather
func (g *Gather[T]) Close() <-chan struct{} {
	close(g.ch)
	return g.done
}

// Put add element
func (g *Gather[T]) Put(ele T) {
	g.ch <- ele
}

// callback call user function
func (g *Gather[T]) callback() {
	if len(g.batch) < 1 {
		return
	}
	batch := g.batch
	g.batch = make([]T, 0, g.size)
	// sync call
	g.call(batch, g)
}
