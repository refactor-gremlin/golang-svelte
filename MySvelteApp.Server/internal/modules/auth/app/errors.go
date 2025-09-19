package app

import "errors"

// ValidationError indicates the payload failed validation rules.
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// ConflictError indicates the request conflicts with existing state (e.g. duplicate username).
type ConflictError struct {
	Message string
}

func (e ConflictError) Error() string {
	return e.Message
}

// UnauthorizedError indicates credentials were invalid.
type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

// IsValidationError returns true when err is a ValidationError.
func IsValidationError(err error) bool {
	var target ValidationError
	return errors.As(err, &target)
}

// IsConflictError returns true when err is a ConflictError.
func IsConflictError(err error) bool {
	var target ConflictError
	return errors.As(err, &target)
}

// IsUnauthorizedError returns true when err is an UnauthorizedError.
func IsUnauthorizedError(err error) bool {
	var target UnauthorizedError
	return errors.As(err, &target)
}
