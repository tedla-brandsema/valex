package valex

import (
	"cmp"
	"errors"
	"fmt"
)

// Validator defines the behavior for validating a value of type T.
type Validator[T any] interface {
	Validate(val T) (ok bool, err error)
}

// ValidatorFunc adapts a function to the Validator interface.
type ValidatorFunc[T any] func(val T) (ok bool, err error)

// Validate calls the underlying function.
func (p ValidatorFunc[T]) Validate(val T) (ok bool, err error) {
	return p(val)
}

// ValidatedValue stores a value and validates updates with the provided Validator.
type ValidatedValue[T cmp.Ordered] struct {
	value     T
	Validator Validator[T]
}

// Set validates and stores the value.
func (v *ValidatedValue[T]) Set(val T) error {
	if v.Validator == nil {
		return errors.New("no validator set")
	}
	if ok, err := v.Validator.Validate(val); !ok {
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
	if ok, err := v.Validate(val); !ok {
		panic(err)
	}
	return val
}
