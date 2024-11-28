package redis

import (
	"encoding/gob"
	"github.com/sasha-s/go-deadlock"
	"ledis-server/utils"
	"os"
	"time"
)

type redis struct {
	data          map[string]Item
	expirationKey map[string]time.Time
	ttl           map[string]time.Duration
	mu            *deadlock.RWMutex
}

func NewRedis() *redis {
	return &redis{
		data:          make(map[string]Item),
		expirationKey: make(map[string]time.Time),
		ttl:           make(map[string]time.Duration),
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
	return s.getOrExpired(key)
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

func (s *redis) FlushDB() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]Item)
	s.expirationKey = make(map[string]time.Time)
	s.ttl = make(map[string]time.Duration)
}

func (s *redis) Expire(key string, ttlInSeconds int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if item, exist := s.getOrExpired(key); !exist || item == nil {
		return utils.ErrKeyDoesNotExist(key)
	}

	ttl := time.Duration(ttlInSeconds) * time.Second
	s.expirationKey[key] = time.Now().Add(ttl)
	s.ttl[key] = ttl

	return nil
}

func (s *redis) TTL(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if item, exist := s.getOrExpired(key); !exist || item == nil {
		return -1, utils.ErrKeyDoesNotExist(key)
	}

	ttl, exists := s.ttl[key]
	if !exists {
		return -1, nil
	}
	return int(ttl / 1e9), nil
}

func (s *redis) Gets(keys ...string) []Item {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := make([]Item, 0, len(keys))
	for _, key := range keys {
		item, exists := s.getOrExpired(key)
		if exists {
			items = append(items, item)
		}
	}

	return items
}

func (s *redis) MakeSnapshot() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tempFileName := "snapshot_temp.rdb"
	finalFileName := "snapshot.rdb"

	snapshot := struct {
		Data          map[string]Item
		ExpirationKey map[string]time.Time
		TTL           map[string]time.Duration
	}{
		Data:          s.data,
		ExpirationKey: s.expirationKey,
		TTL:           s.ttl,
	}

	file, err := os.Create(tempFileName)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		if err != nil {
			os.Remove(tempFileName)
		}
	}()

	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(snapshot); err != nil {
		return err
	}

	if err = os.Rename(tempFileName, finalFileName); err != nil {
		return err
	}

	return nil
}

func (s *redis) LoadSnapshot() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open("snapshot.rdb")
	if err != nil {
		return err
	}
	defer file.Close()

	snapshot := struct {
		Data          map[string]Item
		ExpirationKey map[string]time.Time
		TTL           map[string]time.Duration
	}{}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&snapshot)
	if err != nil {
		return err
	}

	s.data = snapshot.Data
	s.expirationKey = snapshot.ExpirationKey
	s.ttl = snapshot.TTL

	return nil
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
	delete(s.ttl, key)
}

func (s *redis) getOrExpired(key string) (Item, bool) {
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
