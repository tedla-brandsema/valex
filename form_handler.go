package valex

import (
	"errors"
	"net/http"

	"github.com/tedla-brandsema/tagex"
)

// FormError wraps a validation error with an HTTP status code.
type FormError struct {
	Status int
	Err    error
}

func (e *FormError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *FormError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// StatusCode returns the associated HTTP status code.
func (e *FormError) StatusCode() int {
	if e == nil {
		return 0
	}
	return e.Status
}

// FormStatus maps validation errors to HTTP status codes.
func FormStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var tagErr *tagex.TagError
	if errors.As(err, &tagErr) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrFieldRequired) {
		return http.StatusUnprocessableEntity
	}
	return http.StatusBadRequest
}

// ValidateForm wraps NewFormValidator + Validate and returns a FormError on failure.
func ValidateForm(r *http.Request, dst any) (bool, error) {
	validator, err := NewFormValidator(r)
	if err != nil {
		return false, &FormError{Status: http.StatusBadRequest, Err: err}
	}
	ok, err := validator.Validate(dst)
	if err != nil {
		return false, &FormError{Status: FormStatus(err), Err: err}
	}
	return ok, nil
}
