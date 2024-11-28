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
	Set(key string, value Item)
	Get(key string) (Item, bool)
}
