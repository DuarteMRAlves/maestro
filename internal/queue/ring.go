package queue

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"sync"
)

type Ring interface {
	Push(v interface{})
	Pop() interface{}
	Len() int
}

// ring is a simple implementation of the Ring interface. It uses a global lock
// hindering its performance. If performance ever becomes a bottleneck, a custom
// implementation with buffers will have to be implemented.
type ring struct {
	buf  chan interface{}
	cond *sync.Cond
}

func NewRing(capacity int) (Ring, error) {
	if capacity <= 0 {
		return nil, errdefs.InvalidArgumentWithMsg("capacity must be positive")
	}
	r := &ring{
		buf:  make(chan interface{}, capacity),
		cond: sync.NewCond(&sync.Mutex{}),
	}
	return r, nil
}

func (r *ring) Push(v interface{}) {
	r.cond.L.Lock()
	defer r.cond.L.Unlock()
	// If full discard one element
	if len(r.buf) == cap(r.buf) {
		<-r.buf
	}
	r.buf <- v
	r.cond.Signal()
}

func (r *ring) Pop() interface{} {
	r.cond.L.Lock()
	defer r.cond.L.Unlock()
	for len(r.buf) == 0 {
		r.cond.Wait()
	}
	return <-r.buf
}

func (r *ring) Len() int {
	return len(r.buf)
}
