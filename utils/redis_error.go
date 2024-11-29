package utils

import (
	"errors"
	"fmt"
)

var (
	ErrArgsLengthNotMatch = NewErrorResponse(
		errors.New("args length is not match"), "args length is not match", "args length is not match",
		"ErrArgsLengthIsNotMatch",
	)
	ErrCommandDoesNotExist = NewErrorResponse(
		errors.New("this command does not exist"), "this command does not exist", "this command does not exist",
		"ErrCommandDoesNotExist",
	)
)

func ErrCommandRegisteredDuplicate(command string) error {
	err := fmt.Errorf("command %s already exists", command)
	msg := fmt.Sprintf("command %s already exists", command)
	return NewErrorResponse(err, msg, msg, "ErrCommandRegisteredDuplicate")
}

func ErrKeyDoesNotExist(key string) error {
	err := fmt.Errorf("key %s does not exist", key)
	msg := fmt.Sprintf("key %s does not exist", key)
	return NewErrorResponse(err, msg, msg, "ErrKeyDoesNotExist")
}

func ErrTypeMismatch(wantTypeId, currentTypeId int) error {
	err := fmt.Errorf(
		"type mismatch, this key is not a %s, it is a %s", TypeToString[wantTypeId], TypeToString[currentTypeId],
	)
	msg := fmt.Sprintf(
		"type mismatch, this key is not a %s, it is a %s", TypeToString[wantTypeId], TypeToString[currentTypeId],
	)
	return NewErrorResponse(err, msg, msg, "ErrTypeMismatch")
}
