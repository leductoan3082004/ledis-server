package redis

import (
	"encoding/gob"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockItem struct {
	Val any
	Typ int
}

func (m *MockItem) Value() any {
	return m.Val
}

func (m *MockItem) Type() int {
	return m.Typ
}

func TestRedis_SetAndGet(t *testing.T) {
	redis := NewRedis()

	item := &MockItem{Val: "test Val", Typ: 1}

	redis.Set("key1", item)

	retrievedItem, exists := redis.GetOrExpired("key1")
	assert.True(t, exists, "Item should exist in Redis")
	assert.Equal(t, item.Value(), retrievedItem.Value(), "The retrieved item's Val should match the original")
}

func TestRedis_Delete(t *testing.T) {
	redis := NewRedis()

	item := &MockItem{Val: "test Val", Typ: 1}
	redis.Set("key2", item)

	redis.Delete("key2")

	_, exists := redis.GetOrExpired("key2")
	assert.False(t, exists, "Item should be deleted")
}

func TestRedis_Expire(t *testing.T) {
	redis := NewRedis()

	item := &MockItem{Val: "test Val", Typ: 1}
	redis.Set("key3", item)

	err := redis.Expire("key3", 2)
	assert.NoError(t, err, "Setting expiration should not return an error")

	ttl, err := redis.TTL("key3")
	assert.NoError(t, err, "Fetching TTL should not return an error")
	assert.Greater(t, ttl, 0, "TTL should be greater than 0")

	time.Sleep(3 * time.Second)

	_, exists := redis.GetOrExpired("key3")
	assert.False(t, exists, "The item should have expired")
}

func TestRedis_Keys(t *testing.T) {
	redis := NewRedis()

	item1 := &MockItem{Val: "value1", Typ: 1}
	item2 := &MockItem{Val: "value2", Typ: 1}
	redis.Set("key5", item1)
	redis.Set("key6", item2)

	keys := redis.Keys()
	assert.Contains(t, keys, "key5", "key5 should be in the list of keys")
	assert.Contains(t, keys, "key6", "key6 should be in the list of keys")

	redis.Expire("key5", 1)
	time.Sleep(2 * time.Second)

	keys = redis.Keys()
	assert.NotContains(t, keys, "key5", "key5 should have expired and not be in the list of keys")
	assert.Contains(t, keys, "key6", "key6 should still be in the list of keys")
}

func TestRedis_FlushDB(t *testing.T) {
	redis := NewRedis()

	item1 := &MockItem{Val: "value1", Typ: 1}
	item2 := &MockItem{Val: "value2", Typ: 1}
	redis.Set("key7", item1)
	redis.Set("key8", item2)

	redis.FlushDB()

	keys := redis.Keys()
	assert.Empty(t, keys, "After flushing, the database should be empty")
}

func TestRedis_MakeAndLoadSnapshot(t *testing.T) {
	gob.RegisterName("MockItem", &MockItem{})
	defer os.Remove("snapshot.rdb")
	redis := NewRedis()

	item1 := &MockItem{Val: "value1", Typ: 1}
	item2 := &MockItem{Val: "value2", Typ: 1}
	redis.Set("key9", item1)
	redis.Set("key10", item2)

	err := redis.MakeSnapshot()
	assert.NoError(t, err, "Making a snapshot should not return an error")

	redis2 := NewRedis()
	err = redis2.LoadSnapshot()
	assert.NoError(t, err, "Loading a snapshot should not return an error")

	retrievedItem1, exists1 := redis2.Get("key9")
	retrievedItem2, exists2 := redis2.Get("key10")
	assert.True(t, exists1, "key9 should exist after loading snapshot")
	assert.True(t, exists2, "key10 should exist after loading snapshot")
	assert.Equal(t, item1.Value(), retrievedItem1.Value(), "The loaded item for key9 should match")
	assert.Equal(t, item2.Value(), retrievedItem2.Value(), "The loaded item for key10 should match")
}

func TestRedis_ExpiredKey(t *testing.T) {
	redis := NewRedis()

	item := &MockItem{Val: "test Val", Typ: 1}
	redis.Set("key11", item)

	err := redis.Expire("key11", 2)
	assert.NoError(t, err, "Setting expiration should not return an error")

	assert.False(t, redis.Expired("key11"), "The key should not be expired immediately")

	time.Sleep(3 * time.Second)

	assert.True(t, redis.Expired("key11"), "The key should be expired after the TTL")
}

func TestRedis_GetOrExpired(t *testing.T) {
	redis := NewRedis()

	item := &MockItem{Val: "test Val", Typ: 1}
	redis.Set("key12", item)

	redis.Expire("key12", 1)
	time.Sleep(2 * time.Second)

	retrievedItem, exists := redis.GetOrExpired("key12")
	assert.False(t, exists, "The key should have expired and should not be retrieved")
	assert.Nil(t, retrievedItem, "The retrieved item should be nil")
}

func TestRedis_ExpirePeriodically(t *testing.T) {
	redis := NewRedis()
	item := &MockItem{Val: "test Val", Typ: 1}
	redis.Set("key13", item)
	redis.Expire("key13", 1)

	// sleep 6 seconds to make sure active expiration can run
	time.Sleep(6 * time.Second)

	redis.mu.Lock()
	_, exist := redis.data["key13"]
	redis.mu.Unlock()

	assert.False(t, exist, "The key should have expired and should not be retrieved")
}
