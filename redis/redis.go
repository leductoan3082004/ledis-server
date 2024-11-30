package redis

import (
	"encoding/gob"
	"github.com/sasha-s/go-deadlock"
	"ledis-server/logging"
	"ledis-server/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type redis struct {
	data          map[string]Item
	expirationKey map[string]time.Time
	mu            *deadlock.RWMutex
	signalChan    chan os.Signal
	done          chan struct{}
	stopped       bool
}

func NewRedis() *redis {
	rds := &redis{
		data:          make(map[string]Item),
		expirationKey: make(map[string]time.Time),
		mu:            new(deadlock.RWMutex),
		signalChan:    make(chan os.Signal, 1),
		done:          make(chan struct{}),
	}

	rds.expireKeysPeriodically()
	signal.Notify(rds.signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		select {
		case <-rds.signalChan:
			logging.GetLogger().Info("Received termination signal. Shutting down...")
			close(rds.done)
			rds.Stop()
		}
	}()

	return rds
}

func (s *redis) expireKeysPeriodically() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case <-s.done:
				logging.GetLogger().Info("ExpirationKeysPeriodically gracefully exit")
				return

			case <-ticker.C:
				s.mu.Lock()
				for key := range s.expirationKey {
					s.getOrExpired(key)
				}
				s.mu.Unlock()
			}
		}
	}()

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

func (s *redis) GetOrExpired(key string) (Item, bool) {
	return s.getOrExpired(key)
}

func (s *redis) Delete(key string) {
	s.delete(key)
}

func (s *redis) Expired(key string) bool {
	return s.expired(key)
}

func (s *redis) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		if !s.expired(k) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (s *redis) FlushDB() {
	s.data = make(map[string]Item)
	s.expirationKey = make(map[string]time.Time)
}

func (s *redis) Expire(key string, ttlInSeconds int) error {
	if item, exist := s.getOrExpired(key); !exist || item == nil {
		return utils.ErrKeyDoesNotExist(key)
	}

	ttl := time.Duration(ttlInSeconds) * time.Second
	s.expirationKey[key] = time.Now().Add(ttl)

	return nil
}

func (s *redis) TTL(key string) (int, error) {
	_, exist := s.getOrExpired(key)
	if !exist {
		return -1, utils.ErrKeyDoesNotExist(key)
	}
	expirationTime, exists := s.expirationKey[key]
	if !exists {
		return -1, nil
	}
	ttl := expirationTime.Sub(time.Now())
	return int(ttl / 1e9), nil
}

func (s *redis) Gets(keys ...string) []Item {
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
	tempFileName := "snapshot_temp.rdb"
	finalFileName := "snapshot.rdb"

	snapshot := struct {
		Data          map[string]Item
		ExpirationKey map[string]time.Time
		TTL           map[string]time.Duration
	}{
		Data:          s.data,
		ExpirationKey: s.expirationKey,
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

	return nil
}

func (s *redis) Stop() error {
	<-s.done

	s.mu.Lock()
	s.stopped = true
	s.mu.Unlock()

	logging.GetLogger().Info("Shutting down...")
	logging.GetLogger().Info("Shutdown Redis successfully")

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
