package util

import (
	"fmt"
	"sync"
)

type Pool[T any] struct {
	mu      *sync.Mutex
	items   []T
	creator func() (T, error)
}

func NewPool[T any](creator func() (T, error)) *Pool[T] {
	return &Pool[T]{
		creator: creator,
		mu:      &sync.Mutex{},
	}
}

func (p *Pool[T]) Get() (T, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.items) == 0 {
		item, err := p.creator()
		if err != nil {
			return item, fmt.Errorf("create pool item: %w", err)
		}
		p.items = append(p.items, item)
	}

	result := p.items[0]
	p.items = p.items[1:]
	return result, nil
}

func (p *Pool[T]) Put(item T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.items = append(p.items, item)
}
