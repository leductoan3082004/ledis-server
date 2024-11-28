package types

import (
	"container/list"
	"ledis-server/redis"
	"ledis-server/utils"
)

type ListType struct {
	list *list.List
}

func NewListType() redis.Item {
	return &ListType{
		list: list.New(),
	}
}

func (s *ListType) Value() any {
	return s.list
}

func (s *ListType) Type() int {
	return utils.ListType
}

func (s *ListType) LLen() int {
	return s.list.Len()
}

func (s *ListType) LPush(values ...*string) int {
	for _, v := range values {
		s.list.PushFront(*v)
	}
	return s.LLen()
}

func (s *ListType) RPush(values ...*string) int {
	for _, v := range values {
		s.list.PushBack(*v)
	}
	return s.LLen()
}

func (s *ListType) LPop() *string {
	if value := s.list.Front(); value != nil {
		s.list.Remove(value)
		return utils.PtrToValue[string](value)
	}
	return nil
}

func (s *ListType) RPop() *string {
	if value := s.list.Back(); value != nil {
		s.list.Remove(value)
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

	e := atIndex(from, s.list)
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
