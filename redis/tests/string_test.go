package tests

import (
	"github.com/stretchr/testify/assert"
	"ledis-server/utils"
	"testing"
	"time"
)

func TestStringCommands(t *testing.T) {
	t.Run(
		"Test simple string commands", func(t *testing.T) {
			cm := initDB()
			_, err := cm.Execute("SET", "key1", "value1")
			assert.NoError(t, err)

			res, err := cm.Execute("GET", "key1")
			assert.NoError(t, err)
			assert.Equal(t, "value1", res)

			// this key does not have any TTL so it will be -1
			ttl, err := cm.Execute("TTL", "key1")
			assert.NoError(t, err)
			assert.Equal(t, -1, ttl)

			// set new ttl for this key, the value is 2 seconds
			_, err = cm.Execute("EXPIRE", "key1", "2")
			assert.NoError(t, err)

			// now get again the ttl, this will be 2
			ttl, err = cm.Execute("TTL", "key1")
			assert.NoError(t, err)
			assert.Equal(t, 2, ttl)

			// sleep 2 seconds to make it expired and get again, it will return error not found
			time.Sleep(2 * time.Second)
			ttl, err = cm.Execute("TTL", "key1")
			assert.Error(t, err)
			assert.True(t, assert.ObjectsAreEqualValues(err, utils.ErrKeyDoesNotExist("key1")))
		},
	)

}
