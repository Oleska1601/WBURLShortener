package errs

import (
	"errors"
	"fmt"
)

var (
	NotFoundError      = errors.New("not found")
	AlreadyExistsError = errors.New("already exists")
	ConflictError      = errors.New("conflict")
)

func NewNotFoundError(msg string) error {
	return fmt.Errorf("%w: %s", NotFoundError, msg)
}

func NewAlreadyExistsError(msg string) error {
	return fmt.Errorf("%w: %s", AlreadyExistsError, msg)
}

func NewConflictError(msg string) error {
	return fmt.Errorf("%w: %s", ConflictError, msg)
}
