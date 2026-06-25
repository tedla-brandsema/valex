package forms

import (
	"errors"
	"net/http"

	"github.com/tedla-brandsema/valex"
)

// Error wraps a validation error with an HTTP status code.
type Error struct {
	Status int
	Err    error
}

// Error returns the wrapped error's message, or a generic message when no inner
// error is set.
func (e *Error) Error() string {
	if e == nil || e.Err == nil {
		return "forms: validation failed"
	}
	return e.Err.Error()
}

// Unwrap returns the wrapped error so that errors.Is and errors.As can inspect
// the underlying validation or binding failure.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// StatusCode returns the associated HTTP status code.
func (e *Error) StatusCode() int {
	if e == nil {
		return 0
	}
	return e.Status
}

// Status maps validation errors to HTTP status codes.
func Status(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var tagErr *valex.TagError
	if errors.As(err, &tagErr) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrFieldRequired) {
		return http.StatusUnprocessableEntity
	}
	return http.StatusBadRequest
}

// Validate parses the request, binds and validates dst, and returns an *Error
// (with an HTTP status code) on failure.
func Validate(r *http.Request, dst any) (bool, error) {
	validator, err := New(r)
	if err != nil {
		return false, &Error{Status: http.StatusBadRequest, Err: err}
	}
	return validator.Validate(dst)
}
