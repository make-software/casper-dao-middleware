package errors

import (
	"net/http"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

const DuplicateEntryErrorCode = 1062

type DuplicateEntityError struct {
	BaseError
}

func NewDuplicateEntityError(message string) error {
	return &DuplicateEntityError{
		BaseError{
			Code:     "duplicate_entity",
			Message:  message,
			HTTPCode: http.StatusConflict,
		},
	}
}

func IgnoreDuplicateEntryError(err error) error {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return err
	}
	if me.Number != DuplicateEntryErrorCode {
		return err
	}

	zap.S().With(err).Debug("DuplicateEntry error ignored\n")
	return nil
}

func WrapIfDuplicateEntryError(err error) error {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return err
	}

	if me.Number != DuplicateEntryErrorCode {
		return err
	}

	return NewDuplicateEntityError("duplicate entity")
}
