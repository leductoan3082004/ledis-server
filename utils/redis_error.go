package utils

import (
	"errors"
	"fmt"
)

var (
	ErrArgsLengthNotMatch  = errors.New("args length not match")
	ErrCommandDoesNotExist = errors.New("command does not exist")
)

func ErrCommandRegisteredDuplicate(command string) error {
	return fmt.Errorf("command %s already exists", command)
}

func ErrKeyDoesNotExist(key string) error {
	return fmt.Errorf("key %s does not exist", key)
}
