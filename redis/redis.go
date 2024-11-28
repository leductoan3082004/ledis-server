package redis

import (
	"github.com/sasha-s/go-deadlock"
	"time"
)

type redis struct {
	data          map[string]Item
	expirationKey map[string]time.Time
	mu            *deadlock.RWMutex
}

func NewRedis() *redis {
	return &redis{
		data:          make(map[string]Item),
		expirationKey: make(map[string]time.Time),
		mu:            new(deadlock.RWMutex),
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
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *redis) Get(key string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exist := s.data[key]
	return value, exist
}

func (s *redis) GetOrExpired(key string) (Item, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exist := s.data[key]
	if !exist {
		return nil, false
	}

	if s.expired(key) {
		s.delete(key)
		return nil, false
	}
	return value, true
}

func (s *redis) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.delete(key)
}

func (s *redis) Expired(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.expired(key)
}

func (s *redis) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		if !s.expired(k) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (s *redis) expired(key string) bool {
	return s.hasTTLSet(key) && s.keyHasExpired(key)
}

func (s *redis) hasTTLSet(key string) bool {
	_, ok := s.expirationKey[key]
	return ok
}

func (s *redis) keyHasExpired(key string) bool {
	return time.Now().After(s.expirationKey[key])
}

func (s *redis) delete(key string) {
	delete(s.data, key)
	delete(s.expirationKey, key)
}
