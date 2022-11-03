package casper

import (
	"strings"
)

type BaseError struct {
	message string
}

func (e *BaseError) Error() string {
	return e.message
}

type ValueNotFoundError struct {
	BaseError
}

func NewValueNotFoundError(message string) *ValueNotFoundError {
	return &ValueNotFoundError{
		BaseError{
			message: message,
		},
	}
}

type RootNotFoundError struct {
	BaseError
}

func NewRootNotFoundError(message string) *RootNotFoundError {
	return &RootNotFoundError{
		BaseError{
			message: message,
		},
	}
}

func newErrorFromRPCError(err error) error {
	switch {
	case strings.Contains(err.Error(), "ValueNotFound"):
		valNotFoundErr := NewValueNotFoundError(err.Error())

		return valNotFoundErr
	case strings.Contains(err.Error(), "RootNotFound"):
		rootNotFoundErr := NewRootNotFoundError(err.Error())

		return rootNotFoundErr
	default:
		return err
	}
}
