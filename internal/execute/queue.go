package execute

import "sync"

type queue struct {
	data chan state
	mu   sync.Mutex
}

func newQueue(cap int) *queue {
	var q queue
	q.data = make(chan state, cap)
	return &q
}

func (q *queue) push(s state) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) == cap(q.data) {
		<-q.data
	}
	q.data <- s
}

func (q *queue) pop() state {
	q.mu.Lock()
	defer q.mu.Unlock()
	return <-q.data
}
