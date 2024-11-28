package types

import (
	"ledis-server/redis"
	"ledis-server/utils"
)

type StringType struct {
	Val *string
}

func NewStringType(val string) redis.Item {
	return &StringType{
		Val: &val,
	}
}

func (s *StringType) Value() any {
	return s.Val
}

func (s *StringType) Type() int {
	return utils.StringType
}
