package tests

import (
	"github.com/stretchr/testify/assert"
	"ledis-server/utils"
	"testing"
)

func TestSADD(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("SADD", "set1", "value1", "value2", "value3")
	assert.NoError(t, err)

	members, err := cm.Execute("SMEMBERS", "set1")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"value1", "value2", "value3"}, members)
}

func TestSCARD(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("SADD", "set1", "value1", "value2", "value3")
	assert.NoError(t, err)

	cardinality, err := cm.Execute("SCARD", "set1")
	assert.NoError(t, err)
	assert.Equal(t, 3, cardinality)

	_, err = cm.Execute("DEL", "set1")
	assert.NoError(t, err)

	cardinality, err = cm.Execute("SCARD", "set1")
	assert.True(t, assert.ObjectsAreEqualValues(err, utils.ErrKeyDoesNotExist("set1")))
	assert.Equal(t, nil, cardinality)
}

func TestSMEMBERS(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("SADD", "set1", "value1", "value2", "value3")
	assert.NoError(t, err)

	members, err := cm.Execute("SMEMBERS", "set1")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"value1", "value2", "value3"}, members)
}

func TestSREM(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("SADD", "set1", "value1", "value2", "value3")
	assert.NoError(t, err)

	_, err = cm.Execute("SREM", "set1", "value1", "value2")
	assert.NoError(t, err)

	members, err := cm.Execute("SMEMBERS", "set1")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"value3"}, members)
}

func TestSINTER(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("SADD", "set1", "value1", "value2", "value3")
	assert.NoError(t, err)
	_, err = cm.Execute("SADD", "set2", "value2", "value3", "value4")
	assert.NoError(t, err)

	intersection, err := cm.Execute("SINTER", "set1", "set2")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"value2", "value3"}, intersection)

	_, err = cm.Execute("SADD", "set3", "value5")
	assert.NoError(t, err)
	intersection, err = cm.Execute("SINTER", "set1", "set3")
	assert.NoError(t, err)
	assert.Empty(t, intersection)
}

func TestSADDWithDuplicates(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("SADD", "set1", "value1", "value1", "value2")
	assert.NoError(t, err)

	members, err := cm.Execute("SMEMBERS", "set1")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"value1", "value2"}, members)
}

func TestSADDWithEmptySet(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("SADD", "set1", "value1", "value2")
	assert.NoError(t, err)

	members, err := cm.Execute("SMEMBERS", "set1")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"value1", "value2"}, members)
}

func TestSREMKeyDoesNotExist(t *testing.T) {
	cm := initDB()

	_, err := cm.Execute("SREM", "nonexistentSet", "value1")
	assert.True(t, assert.ObjectsAreEqualValues(err, utils.ErrKeyDoesNotExist("nonexistentSet")))
}
