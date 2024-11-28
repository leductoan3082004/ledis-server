package redis

type Item interface {
	Value() any
	Type() int
}

type Redis interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()

	// Set key equal to value
	Set(key string, value Item)

	// Get just get but not check expiration
	Get(key string) (Item, bool)

	// Expired check if key has expired or not
	Expired(key string) bool

	// GetAndExpired lazy key expiration when get it
	GetAndExpired(key string) (Item, bool)
}
