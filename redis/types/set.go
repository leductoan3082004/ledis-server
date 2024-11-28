package types

import (
	"github.com/scylladb/go-set/strset"
	"ledis-server/redis"
	"ledis-server/utils"
)

type SetType struct {
	set *strset.Set
}

func NewSetType() redis.Item {
	return &SetType{
		set: strset.New(),
	}
}

// Value Return the underlying set as a value
func (s *SetType) Value() any {
	return s.set
}

// Type Return the type identifier for SetType
func (s *SetType) Type() int {
	return utils.SetType
}

// SAdd SADD: Add values to the set
func (s *SetType) SAdd(values ...string) int {
	initialSize := s.set.Size()
	s.set.Add(values...)
	return s.set.Size() - initialSize
}

func (s *SetType) GetSet() *strset.Set {
	return s.set
}

// SCard Return the number of elements in the set
func (s *SetType) SCard() int {
	return s.set.Size()
}

// SMembers Return all members of the set as a slice
func (s *SetType) SMembers() []string {
	return s.set.List()
}

// SRem Remove specified values from the set
func (s *SetType) SRem(values ...string) int {
	initialSize := s.set.Size()
	s.set.Remove(values...)
	return initialSize - s.set.Size()
}

// SInter Perform set intersection with other sets and return the result as a slice
func SInter(sets ...*strset.Set) []string {
	if len(sets) == 0 {
		return []string{}
	}

	result := strset.Intersection(sets...)

	return result.List()
}
