package queue

import (
	"gotest.tools/v3/assert"
	"sync"
	"testing"
)

func TestRing_NotFull(t *testing.T) {
	capacity := 1000
	r, err := NewRing(capacity)
	assert.NilError(t, err, "ring creation")

	r.Push(1)
	pop, ok := r.Pop().(int)
	assert.Assert(t, ok, "pop correct type")
	assert.Equal(t, 1, pop, "pop correct value")

	assert.Equal(t, 0, r.Len(), "empty ring")

	for i := 0; i < capacity; i++ {
		r.Push(i)
	}
	for i := 0; i < capacity; i++ {
		pop, ok = r.Pop().(int)
		assert.Assert(t, ok, "pop cycle correct type")
		assert.Equal(t, i, pop, "pop cycle correct value")
	}

	assert.Equal(t, 0, r.Len(), "empty ring after cycle")
}

func TestRing_Full(t *testing.T) {
	capacity := 1000
	r, err := NewRing(capacity)
	assert.NilError(t, err, "ring creation")

	for i := 0; i < 2*capacity; i++ {
		r.Push(i)
	}

	assert.Equal(t, capacity, r.Len(), "full ring after push cycle")

	// Only the elements between capacity and 2*capacity should be stored
	for i := capacity; i < 2*capacity; i++ {
		pop, ok := r.Pop().(int)
		assert.Assert(t, ok, "pop correct type")
		assert.Equal(t, i, pop, "pop correct value")
	}

	assert.Equal(t, 0, r.Len(), "empty ring after pop cycle")
}

func TestRing_ConcurrentModification(t *testing.T) {
	var wg sync.WaitGroup

	capacity := 1000

	r, err := NewRing(capacity)
	assert.NilError(t, err, "ring creation")

	wg.Add(1)

	startChan := make(chan struct{}, 1)

	go func() {
		startChan <- struct{}{}
		for i := 0; i < capacity; i++ {
			pop, ok := r.Pop().(int)
			assert.Assert(t, ok, "pop cycle correct type")
			assert.Equal(t, i, pop, "pop cycle correct value")
		}
		wg.Done()
	}()

	<-startChan
	for i := 0; i < capacity; i++ {
		r.Push(i)
	}

	wg.Wait()

	assert.Equal(t, 0, r.Len(), "empty ring after cycle")
}
