package common

import (
	"errors"
	"fmt"
)

// "Root cause" categories
var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrAuthFailed       = errors.New("auth failed")
	ErrAccessDenied     = errors.New("access denied")
	ErrValidation       = errors.New("validation failed")
	ErrConflict         = errors.New("conflict")
	ErrDBOperation      = errors.New("db operation failed")
)

// AppError wraps both a generic error and a categorized error
type AppError struct {
	Generic error // from library/DB/etc
	Kind    error // mapped category
}

func (e *AppError) Error() string {

	if e.Generic != nil {
		return fmt.Sprintf("%s: %v", e.Kind.Error(), e.Generic)
	}

	return e.Kind.Error()
}

// Unwrap lets errors.Is / errors.As work
func (e *AppError) Unwrap() error {
	return e.Kind
}

func NewError(generic, kind error) error {
	return &AppError{Generic: generic, Kind: kind}
}
