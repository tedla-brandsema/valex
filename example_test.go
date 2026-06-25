package valex_test

import (
	"errors"
	"fmt"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/valex/validators"
)

// ValidatorFunc adapts a plain function into a Validator. A nil return means the
// value is valid.
func ExampleValidatorFunc() {
	nonEmpty := valex.ValidatorFunc[string](func(s string) error {
		if s == "" {
			return errors.New("must not be empty")
		}
		return nil
	})

	fmt.Println(nonEmpty.Validate("hello"))
	fmt.Println(nonEmpty.Validate(""))
	// Output:
	// <nil>
	// must not be empty
}

// ValidatedValue stores a value only when it passes the configured Validator,
// leaving the previous value in place on failure.
func ExampleValidatedValue() {
	positive := valex.ValidatorFunc[int](func(n int) error {
		if n <= 0 {
			return errors.New("must be positive")
		}
		return nil
	})

	v := valex.ValidatedValue[int]{Validator: positive}

	if err := v.Set(42); err == nil {
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

// ValidateStruct validates fields via the "val" tag. Directives are opt-in:
// register the ones you use (here from the valex/validators catalog) first.
func ExampleValidateStruct() {
	valex.RegisterDirective(&validators.EmailValidator{})
	valex.RegisterDirective(&validators.IntRangeValidator{})

	type User struct {
		Email string `val:"email"`
		Age   int    `val:"rangeint,min=0,max=120"`
	}

	fmt.Println(valex.ValidateStruct(&User{Email: "gopher@example.com", Age: 30}))
	fmt.Println(valex.ValidateStruct(&User{Email: "gopher@example.com", Age: 200}))
	// Output:
	// <nil>
	// tag "val" error: directive processing field "Age" directive "rangeint": value 200 is out of range [0, 120]
}

// evenDirective is a custom "val" directive that accepts only even ints.
// A directive is registered as a pointer so tagex can populate its parameters.
type evenDirective struct{}

func (*evenDirective) Name() string              { return "even" }
func (*evenDirective) Mode() tagex.DirectiveMode { return tagex.EvalMode }
func (*evenDirective) Handle(n int) (int, error) {
	if n%2 != 0 {
		return n, fmt.Errorf("value %d is not even", n)
	}
	return n, nil
}

// RegisterDirective extends the "val" tag with a custom directive.
func ExampleRegisterDirective() {
	valex.RegisterDirective(&evenDirective{})

	type Ticket struct {
		Seats int `val:"even"`
	}

	fmt.Println(valex.ValidateStruct(&Ticket{Seats: 4}))
	fmt.Println(valex.ValidateStruct(&Ticket{Seats: 3}))
	// Output:
	// <nil>
	// tag "val" error: directive processing field "Seats" directive "even": value 3 is not even
}
