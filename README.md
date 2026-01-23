# Valex

*Valex* is an extensible validation library for Go. It provides a flexible way to define type-safe validators and wrap 
values in a way that ensures they satisfy custom validation rules before being set. 

Features

* **Generic Validators:** Define validators for any ordered type (e.g. integers, floats, strings).
* **Validator Interface & Adapter:** Implement your own validation logic via the `Validator[T]` interface or create quick validators using the `ValidatorFunc[T]` adapter.
* **Validated Value Wrapper:** Use the `ValidatedValue[T]` type to ensure that only valid values (as determined by your validator) are set.
* **Tag-Based Validation:** Use struct tags with built-in directives via `ValidateStruct`.
* **Custom Directives:** Register your own tag directives with `RegisterDirective`.

## Installation

To add Valex to your project, run:

```
go get -u github.com/tedla-brandsema/valex@latest
```

## Examples 

### Defining a Custom Validator

Implement the `Validator[T]` interface for your type. For example, here’s a simple integer range validator:

```go
package main

import (
	"fmt"
	"github.com/tedla-brandsema/valex"
)

type IntRangeValidator struct {
	Min int
	Max int
}

func (v IntRangeValidator) Validate(val int) (bool, error) {
	if val < v.Min || val > v.Max {
		return false, fmt.Errorf("value %d is out of range [%d, %d]", val, v.Min, v.Max)
	}
	return true, nil
}

func main() {
	// Create a Validator
	v := IntRangeValidator{
		Min: 1,
		Max: 10,
	}

	if ok, err := v.Validate(11); !ok {
		fmt.Println("Error:", err)
	}

	// Or use a Validator in conjunction with a ValidatedValue
	vv := valex.ValidatedValue[int]{
		Validator: v,
	}
	
	if err := vv.Set(5); err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	fmt.Println("Validated value:", vv.Get())
}
```

### Using ValidatorFunc

You can also use the `ValidatorFunc[T]` adapter to quickly create validators from functions:

```go
package main

import (
	"fmt"
	"github.com/tedla-brandsema/valex"
)

func main() {
	// Create a validator for strings that ensures they are non-empty.
	nonEmptyValidator := valex.ValidatorFunc[string](func(val string) (bool, error) {
		if val == "" {
			return false, fmt.Errorf("string cannot be empty")
		}
		return true, nil
	})

	vv := valex.ValidatedValue[string]{
		Validator: nonEmptyValidator,
	}

	if err := vv.Set("hello world"); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Validated value:", vv.Get())
}
```

### Tag-Based Validation

You can validate struct fields by using the `val` tag and calling `ValidateStruct`:

```go
package main

import (
	"fmt"
	"github.com/tedla-brandsema/valex"
)

type User struct {
	Name  string `val:"min,size=3"`
	Email string `val:"email"`
	Age   int    `val:"range,min=0,max=120"`
}

func main() {
	user := &User{Name: "Al", Email: "invalid", Age: 200}
	ok, err := valex.ValidateStruct(user)
	fmt.Println(ok, err)
}
```

### Registering Custom Directives

Register your own directives to extend the `val` tag:

```go
package main

import (
	"fmt"
	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/tagex"
)

type EvenDirective struct{}

func (d *EvenDirective) Name() string { return "even" }
func (d *EvenDirective) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}
func (d *EvenDirective) Handle(val int) (int, error) {
	if val%2 != 0 {
		return 0, fmt.Errorf("value %d is not even", val)
	}
	return val, nil
}

func main() {
	valex.RegisterDirective(&EvenDirective{})

	type Item struct {
		Count int `val:"even"`
	}

	item := &Item{Count: 3}
	ok, err := valex.ValidateStruct(item)
	fmt.Println(ok, err)
}
```

### Form Validation

`NewFormValidator` calls `request.ParseForm()`, which parses both POST bodies and
URL query parameters. That means you can use `FormValidator` with GET requests
as well as standard form submissions.

### How It Works

* **Validator Interface:**\
    Define a type that implements the method:`Validate(val T) (ok bool, err error)`
    A successful validation should return `true` (with a `nil` error), whereas a failure should return `false` and an appropriate error message.

* **ValidatedValue:**\
    This type holds a value of type `T` along with an associated `Validator[T]`. 
	* `Set(val T) error`: Uses the validator to ensure that only valid values are stored.
	* `Get() T`: Returns the current value.

## Contributing

Contributions, issues, and feature requests are welcome! Please check the issues page if you’d like to contribute.

## License

This project is licensed under the MIT License – see the [LICENSE](https://github.com/tedla-brandsema/valex/blob/main/LICENSE) file for details.
