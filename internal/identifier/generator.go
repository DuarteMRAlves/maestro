package identifier

import (
	"fmt"
	"sync"
)

type Generator = func() (Id, error)

func GenForSize(idSize int) Generator {
	generators.mu.Lock()
	defer generators.mu.Unlock()
	prev, ok := generators.st[idSize]
	if ok {
		return prev
	}
	gen := newGenerator(idSize)
	generators.st[idSize] = gen
	return gen
}

var generators = struct {
	st map[int]Generator
	mu sync.Mutex
}{st: map[int]Generator{}, mu: sync.Mutex{}}

func newGenerator(idSize int) Generator {
	generated := map[Id]bool{}
	return func() (Id, error) {
		newId, err := Rand(idSize)
		if err != nil {
			return Empty(), fmt.Errorf("generate id: %v", err)
		}
		_, idExists := generated[newId]
		for idExists {
			if newId, err = Rand(idSize); err != nil {
				return Empty(), fmt.Errorf("generate id: %v", err)
			}
			_, idExists = generated[newId]
		}
		generated[newId] = true
		return newId, nil
	}
}
