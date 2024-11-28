package tests

import (
	"github.com/stretchr/testify/assert"
	"ledis-server/utils"
	"testing"
	"time"
)

func TestRPUSH(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("RPUSH", "list1", "value1", "value2")
	assert.NoError(t, err)

	lenResp, err := cm.Execute("LLEN", "list1")
	assert.NoError(t, err)
	assert.Equal(t, 2, lenResp)
}

func TestLPOP(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("RPUSH", "list1", "value1", "value2")
	assert.NoError(t, err)

	lpopResp, err := cm.Execute("LPOP", "list1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", lpopResp)

	lenResp, err := cm.Execute("LLEN", "list1")
	assert.NoError(t, err)
	assert.Equal(t, 1, lenResp)
}

func TestRPOP(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("RPUSH", "list1", "value1", "value2")
	assert.NoError(t, err)

	rpopResp, err := cm.Execute("RPOP", "list1")
	assert.NoError(t, err)
	assert.Equal(t, "value2", rpopResp)

	lenResp, err := cm.Execute("LLEN", "list1")
	assert.NoError(t, err)
	assert.Equal(t, 1, lenResp)
}

func TestLLEN(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("RPUSH", "list1", "value1", "value2", "value3")
	assert.NoError(t, err)

	lenResp, err := cm.Execute("LLEN", "list1")
	assert.NoError(t, err)
	assert.Equal(t, 3, lenResp)

	_, err = cm.Execute("DEL", "list1") // Remove the list
	assert.NoError(t, err)

	lenResp, err = cm.Execute("LLEN", "list1")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), utils.ErrKeyDoesNotExist("list1").Error())
}

func TestLRANGE(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("RPUSH", "list1", "value1", "value2", "value3")
	assert.NoError(t, err)

	rangeResp, err := cm.Execute("LRANGE", "list1", "0", "2")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"value1", "value2", "value3"}, rangeResp)

	// Test LRANGE with range out of bounds
	rangeResp, err = cm.Execute("LRANGE", "list1", "2", "5")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"value3"}, rangeResp)
}

func TestDEL(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("RPUSH", "list1", "value1")
	assert.NoError(t, err)

	_, err = cm.Execute("DEL", "list1")
	assert.NoError(t, err)

	_, err = cm.Execute("LLEN", "list1")
	assert.True(t, assert.ObjectsAreEqualValues(err, utils.ErrKeyDoesNotExist("list1")))
}

func TestFLUSHDB(t *testing.T) {
	cm := initDB()

	// Add multiple lists to DB
	_, err := cm.Execute("RPUSH", "list1", "value1")
	assert.NoError(t, err)
	_, err = cm.Execute("RPUSH", "list2", "value2")
	assert.NoError(t, err)

	// Test FLUSHDB (flush all keys)
	_, err = cm.Execute("FLUSHDB")
	assert.NoError(t, err)

	// Ensure the lists are removed
	lenResp, err := cm.Execute("LLEN", "list1")
	assert.ObjectsAreEqualValues(err, utils.ErrKeyDoesNotExist("list1"))
	assert.Equal(t, nil, lenResp)

	lenResp, err = cm.Execute("LLEN", "list2")
	assert.ObjectsAreEqualValues(err, utils.ErrKeyDoesNotExist("list2"))
	assert.Equal(t, nil, lenResp)
}

func TestEXPIREandTTL(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("RPUSH", "list1", "value1")
	assert.NoError(t, err)

	// Test EXPIRE (set TTL of 1 second)
	_, err = cm.Execute("EXPIRE", "list1", "1")
	assert.NoError(t, err)

	// Test TTL (should return 1)
	ttlResp, err := cm.Execute("TTL", "list1")
	assert.NoError(t, err)
	assert.Equal(t, 1, ttlResp)

	// Wait for expiration and check TTL again
	time.Sleep(2 * time.Second)

	ttlResp, err = cm.Execute("TTL", "list1")
	assert.ObjectsAreEqualValues(err, utils.ErrKeyDoesNotExist("list1"))
	assert.Equal(t, -1, ttlResp) // Key should be expired
}
