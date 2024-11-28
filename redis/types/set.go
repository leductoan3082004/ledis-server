package types

import (
	"bytes"
	"encoding/gob"
	"github.com/scylladb/go-set/strset"
	"ledis-server/redis"
	"ledis-server/utils"
)

type SetType struct {
	Set *strset.Set
}

func NewSetType() redis.Item {
	return &SetType{
		Set: strset.New(),
	}
}
func (s *SetType) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(s.Set.List()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *SetType) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	var items []string
	if err := decoder.Decode(&items); err != nil {
		return err
	}

	s.Set = strset.New(items...)
	return nil
}

// Value Return the underlying Set as a Val
func (s *SetType) Value() any {
	return s.Set
}

// Type Return the type identifier for SetType
func (s *SetType) Type() int {
	return utils.SetType
}

// SAdd SADD: Add values to the Set
func (s *SetType) SAdd(values ...string) int {
	initialSize := s.Set.Size()
	s.Set.Add(values...)
	return s.Set.Size() - initialSize
}

func (s *SetType) GetSet() *strset.Set {
	return s.Set
}

// SCard Return the number of elements in the Set
func (s *SetType) SCard() int {
	return s.Set.Size()
}

// SMembers Return all members of the Set as a slice
func (s *SetType) SMembers() []string {
	return s.Set.List()
}

// SRem Remove specified values from the Set
func (s *SetType) SRem(values ...string) int {
	initialSize := s.Set.Size()
	s.Set.Remove(values...)
	return initialSize - s.Set.Size()
}

// SInter Perform Set intersection with other sets and return the result as a slice
func SInter(sets ...*strset.Set) []string {
	if len(sets) == 0 {
		return []string{}
	}

	result := strset.Intersection(sets...)

	return result.List()
}
