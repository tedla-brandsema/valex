// Programmatic validates values in code with ValidatorFunc and ValidatedValue,
// without any struct tags.
package main

import (
	"errors"
	"fmt"

	"github.com/tedla-brandsema/valex"
)

func main() {
	// A validator from a plain function: a nil return means valid.
	positive := valex.ValidatorFunc[int](func(n int) error {
		if n <= 0 {
			return errors.New("must be positive")
		}
		return nil
	})

	// ValidatedValue only stores values that pass, leaving the previous value
	// in place on failure.
	v := valex.ValidatedValue[int]{Validator: positive}

	if err := v.Set(42); err != nil {
		fmt.Println("rejected 42:", err)
	} else {
		fmt.Println("stored:", v.Get())
	}

	if err := v.Set(-1); err != nil {
		fmt.Println("rejected -1:", err)
	}
	fmt.Println("still:", v.Get())

	// Output:
	// stored: 42
	// rejected -1: must be positive
	// still: 42
}
