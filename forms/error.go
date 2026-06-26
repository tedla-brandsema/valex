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

// Status maps validation and binding errors to HTTP status codes: 422
// (Unprocessable Entity) for field-level problems — a validation failure
// (*valex.TagError) or a binding failure (*bindError, including a missing
// required field) — and 400 (Bad Request) for everything else, notably a request
// that could not be parsed at all. A field-level error is 422 whether or not a
// neighbor also failed.
func Status(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var tagErr *valex.TagError
	if errors.As(err, &tagErr) {
		return http.StatusUnprocessableEntity
	}
	var bindErr *bindError
	if errors.As(err, &bindErr) {
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

// ValidateAll parses the request, binds and validates dst against the default
// registry, and collects every binding and validation failure instead of
// stopping at the first. It returns nil on success or an *Error; pass the error
// to FieldErrors for a field-keyed map.
func ValidateAll(r *http.Request, dst any) error {
	return ValidateAllWith(r, dst, nil)
}

// ValidateAllWith is like ValidateAll but validates against reg instead of the
// default registry. A nil reg uses the default.
func ValidateAllWith(r *http.Request, dst any, reg *valex.Registry) error {
	validator, err := NewWith(r, reg)
	if err != nil {
		return &Error{status: http.StatusBadRequest, Err: err}
	}
	return validator.ValidateAll(dst)
}

// FieldErrors flattens err — typically from ValidateAll — into a map from struct
// field path to the error for that field, merging binding and validation
// failures into one view. When a field fails both (for example a non-numeric
// value for an int that also has a range rule), the binding error wins: "not a
// number" is the actionable message, where a range complaint about the unset
// zero value is noise.
//
// Keys are struct field paths (e.g. "Email", "Items[2].SKU"), not request keys
// or display names — translate them when rendering. A nil error yields a nil
// map; non-field errors (such as a request that failed to parse) are omitted, so
// keep the returned error authoritative and render this map on top.
func FieldErrors(err error) map[string]error {
	if err == nil {
		return nil
	}
	// Unwrap the *Error envelope to its underlying (joined) error first: walking
	// the envelope directly would let errors.As dive through the join and collapse
	// every field to the first match.
	if fe, ok := err.(*Error); ok {
		if inner := fe.Unwrap(); inner != nil {
			err = inner
		}
	}
	// Validation errors keyed by field path; bind errors aren't *ProcessError, so
	// valex.FieldErrors skips them.
	m := valex.FieldErrors(err)
	if m == nil {
		m = make(map[string]error)
	}
	// Bind errors take precedence on a same-field collision.
	collectBindErrors(err, m)
	if len(m) == 0 {
		return nil
	}
	return m
}

func collectBindErrors(err error, m map[string]error) {
	if j, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range j.Unwrap() {
			collectBindErrors(e, m)
		}
		return
	}
	var be *bindError
	if errors.As(err, &be) {
		m[be.Field] = err
	}
}
