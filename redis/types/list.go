package types

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"ledis-server/redis"
	"ledis-server/utils"
)

type ListType struct {
	List *list.List
}

func NewListType() redis.Item {
	return &ListType{
		List: list.New(),
	}
}

func (l *ListType) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	var items []any
	for e := l.List.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}

	if err := encoder.Encode(items); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (l *ListType) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	var items []any
	if err := decoder.Decode(&items); err != nil {
		return err
	}

	l.List = list.New()
	for _, item := range items {
		l.List.PushBack(item)
	}
	return nil
}

func (s *ListType) Value() any {
	return s.List
}

func (s *ListType) Type() int {
	return utils.ListType
}

func (s *ListType) LLen() int {
	return s.List.Len()
}

func (s *ListType) LPush(values ...*string) int {
	for _, v := range values {
		s.List.PushFront(*v)
	}
	return s.LLen()
}

func (s *ListType) RPush(values ...*string) int {
	for _, v := range values {
		s.List.PushBack(*v)
	}
	return s.LLen()
}

func (s *ListType) LPop() *string {
	if value := s.List.Front(); value != nil {
		s.List.Remove(value)
		return utils.PtrToValue[string](value)
	}
	return nil
}

func (s *ListType) RPop() *string {
	if value := s.List.Back(); value != nil {
		s.List.Remove(value)
		return utils.PtrToValue[string](value)
	}
	return nil
}

func atIndex(index int, list *list.List) *list.Element {
	index = utils.ToPositiveIndex(index, list.Len())
	e, i := list.Front(), 0
	if e == nil {
		return nil
	}
	for ; e.Next() != nil && i < index; i++ {
		if e.Next() == nil {
			return nil
		}
		e = e.Next()
	}
	return e
}

func (s *ListType) LRange(start int, end int) []string {
	values := make([]string, 0)
	from, to := utils.GetPositiveStartEndIndexes(start, end, s.LLen())

	if from > to {
		return values
	}

	e := atIndex(from, s.List)
	if e == nil {
		return values
	}

	values = append(values, *utils.PtrToValue[string](e))
	for i := from; i < to; i++ {
		e = e.Next()
		values = append(values, *utils.PtrToValue[string](e))
	}
	return values
}
