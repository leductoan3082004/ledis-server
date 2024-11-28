package redis

import (
	"sync"
)

type redis struct {
	data map[string]Item
	mu   *sync.RWMutex
}

func NewRedis() *redis {
	return &redis{
		data: make(map[string]Item),
		mu:   new(sync.RWMutex),
	}
}

func (s *redis) Lock() {
	s.mu.Lock()
}

func (s *redis) Unlock() {
	s.mu.Unlock()
}

func (s *redis) RLock() {
	s.mu.RLock()
}

func (s *redis) RUnlock() {
	s.mu.RUnlock()
}

func (s *redis) Set(key string, value Item) {
	s.data[key] = value
}

func (s *redis) Get(key string) (Item, bool) {
	value, exist := s.data[key]
	return value, exist
}
