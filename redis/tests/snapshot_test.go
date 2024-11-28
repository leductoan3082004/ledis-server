package tests

import (
	"encoding/gob"
	"github.com/stretchr/testify/assert"
	"ledis-server/redis"
	"ledis-server/redis/commands"
	"ledis-server/redis/types"
	"os"
	"testing"
)

func initDB() redis.ICommandManager {
	gob.RegisterName("types.ListType", &types.ListType{})
	gob.RegisterName("types.SetType", &types.SetType{})
	gob.RegisterName("types.StringType", &types.StringType{})

	rds := redis.NewRedis()
	commandManager := commands.NewCommandManager(rds)

	return commandManager
}

func TestStringSetGet(t *testing.T) {
	cm := initDB()

	defer os.Remove("snapshot.rdb")

	resp, err := cm.Execute("SET", "value1", "value2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", resp)

	resp, err = cm.Execute("GET", "value1")
	assert.NoError(t, err)
	assert.Equal(t, "value2", resp)
}

func TestListPushRetrieve(t *testing.T) {
	cm := initDB()

	defer os.Remove("snapshot.rdb")

	resp, err := cm.Execute("RPUSH", "list1", "testkey1", "testkey2")
	assert.NoError(t, err)
	assert.Equal(t, 2, resp)

	beforeRestoreList, err := cm.Execute("LRANGE", "list1", "0", "2")
	assert.NoError(t, err)

	_, err = cm.Execute("SNAPSHOT")
	assert.NoError(t, err)

	newCm := initDB()
	_, err = newCm.Execute("RESTORE")
	assert.NoError(t, err)

	afterRestoreList, err := newCm.Execute("LRANGE", "list1", "0", "2")
	assert.NoError(t, err)
	assert.Equal(t, beforeRestoreList, afterRestoreList)
}

func TestSetRetrieve(t *testing.T) {
	cm := initDB()
	defer os.Remove("snapshot.rdb")
	resp, err := cm.Execute("SADD", "key1", "value1", "value2", "value3", "value3")
	assert.NoError(t, err)
	assert.Equal(t, 3, resp)

	_, err = cm.Execute("SNAPSHOT")
	assert.NoError(t, err)

	_, err = cm.Execute("RESTORE")
	assert.NoError(t, err)

	afterCount, err := cm.Execute("SCARD", "key1")
	assert.NoError(t, err)
	assert.Equal(t, 3, afterCount)

	afterRestoreSet, err := cm.Execute("SMEMBERS", "key1")
	assert.NoError(t, err)

	assert.ElementsMatch(t, []string{"value1", "value2", "value3"}, afterRestoreSet)
}
