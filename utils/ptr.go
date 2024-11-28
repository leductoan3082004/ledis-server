package utils

import "container/list"

func PtrToValue[T any](value *list.Element) *T {
	v := value.Value.(T)
	return &v
}
