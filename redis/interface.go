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

	Delete(key string)

	// Expired check if key has expired or not
	Expired(key string) bool

	// GetOrExpired  lazy key expiration when get it
	GetOrExpired(key string) (Item, bool)

	// Keys return list of keys available (not expired or do not have ttl set)
	Keys() []string

	FlushDB()

	Expire(key string, ttlInSeconds int) error
	TTL(key string) (int, error)
	Gets(keys ...string) []Item

	LoadSnapshot() error
	MakeSnapshot() error
}

type ICommandHandler interface {
	CommandName() string
	Execute(args ...string) (any, error)
}

type ICommandManager interface {
	Register(handler ICommandHandler) ICommandManager
	Execute(command string, args ...string) (any, error)
}
