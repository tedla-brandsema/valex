package valex

import (
	"fmt"
)

// Validator defines the behavior for validating a value of type T. A nil error
// means the value is valid; a non-nil error reports why it is invalid.
type Validator[T any] interface {
	Validate(val T) error
}

// ValidatorFunc adapts a function to the Validator interface.
type ValidatorFunc[T any] func(val T) error

// Validate calls the underlying function.
func (p ValidatorFunc[T]) Validate(val T) error {
	return p(val)
}

// ValidatedValue stores a value and validates updates with the provided
// Validator. It is an in-memory guard, not a serialization type: the stored
// value is unexported and a decoder has no way to supply the Validator, so it
// does not round-trip through encoding/json. For serialized or request input,
// validate with ValidateStruct (the "val" tag) instead.
type ValidatedValue[T any] struct {
	value     T
	Validator Validator[T]
}

// Set validates and stores the value.
func (v *ValidatedValue[T]) Set(val T) error {
	if v.Validator == nil {
		return ErrNoValidator
	}
	if err := v.Validator.Validate(val); err != nil {
		return err
	}
	v.value = val

	return nil
}

// Get returns the stored value.
func (v *ValidatedValue[T]) Get() T {
	return v.value
}

// String returns the string representation of the stored value.
func (v *ValidatedValue[T]) String() string {
	return fmt.Sprintf("%v", v.value)
}

// MustValidate validates a value or panics if validation fails.
func MustValidate[T any](val T, v Validator[T]) T {
	if err := v.Validate(val); err != nil {
		panic(err)
	}
	return val
}
