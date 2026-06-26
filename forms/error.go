package forms

import (
	"errors"
	"net/http"

	"github.com/tedla-brandsema/valex"
)

// Error wraps a validation error with an HTTP status code. Read the status with
// StatusCode and the underlying error with Unwrap.
type Error struct {
	status int
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
	return e.status
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

// Validate parses the request, binds and validates dst against valex's default
// registry, and returns nil on success or an *Error (with an HTTP status code)
// on failure.
func Validate(r *http.Request, dst any) error {
	return ValidateWith(r, dst, nil)
}

// ValidateWith is like Validate but validates against reg instead of the default
// registry. A nil reg uses the default. Use it for an isolated directive set.
func ValidateWith(r *http.Request, dst any, reg *valex.Registry) error {
	validator, err := NewWith(r, reg)
	if err != nil {
		return &Error{status: http.StatusBadRequest, Err: err}
	}
	return validator.Validate(dst)
}
