package types

import (
	"ledis-server/redis"
	"ledis-server/utils"
)

type stringType struct {
	value *string
}

func NewStringType(val string) redis.Item {
	return &stringType{
		value: &val,
	}
}

func (s *stringType) Value() any {
	return s.value
}

func (s *stringType) Type() int {
	return utils.StringType
}
