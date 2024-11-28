package types

import (
	"container/list"
	"ledis-server/redis"
	"ledis-server/utils"
)

type listType struct {
	list *list.List
}

func NewListType() redis.Item {
	return &listType{
		list: list.New(),
	}
}

func (s *listType) Value() any {
	return s.list
}

func (s *listType) Type() int {
	return utils.ListType
}

func (s *listType) LLen() int {
	return s.list.Len()
}

func (s *listType) LPush(values ...*string) int {
	for _, v := range values {
		s.list.PushFront(*v)
	}
	return s.LLen()
}

func (s *listType) RPush(values ...*string) int {
	for _, v := range values {
		s.list.PushBack(*v)
	}
	return s.LLen()
}

func toString(value *list.Element) *string {
	v := value.Value.(string)
	return &v
}
func (s *listType) LPop() *string {
	if value := s.list.Front(); value != nil {
		s.list.Remove(value)
		return toString(value)
	}
	return nil
}

func (s *listType) RPop() *string {
	if value := s.list.Back(); value != nil {
		s.list.Remove(value)
		return toString(value)
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

func (l *listType) LRange(start int, end int) []string {
	values := make([]string, 0)
	from, to := utils.GetPositiveStartEndIndexes(start, end, l.LLen())

	if from > to {
		return values
	}

	e := atIndex(from, l.list)
	if e == nil {
		return values
	}

	values = append(values, *toString(e))
	for i := from; i < to; i++ {
		e = e.Next()
		values = append(values, *toString(e))
	}
	return values
}
