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

You can validate struct fields by using the `val` tag and calling `ValidateStruct`.
Optionally pass additional `tagex.Tag` values to process multiple tags in one pass:

```go
package main

import (
	"fmt"
	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/tagex"
)

type User struct {
	Name  string `val:"min,size=3"`
	Email string `val:"email"`
	Age   int    `val:"rangeint,min=0,max=120"`
}

func main() {
	user := &User{Name: "Al", Email: "invalid", Age: 200}
	tag := tagex.NewTag("extra")
	ok, err := valex.ValidateStruct(user, &tag)
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

For convenience, `ValidateForm` wraps request parsing and validation and returns
a `FormError` with an HTTP status code you can use for responses. If you only
need binding (outside of HTTP handlers), `BindFormValues` accepts `url.Values`
directly.

### How It Works

* **Validator Interface:**\
    Define a type that implements the method:`Validate(val T) (ok bool, err error)`
    A successful validation should return `true` (with a `nil` error), whereas a failure should return `false` and an appropriate error message.

* **ValidatedValue:**\
    This type holds a value of type `T` along with an associated `Validator[T]`. 
	* `Set(val T) error`: Uses the validator to ensure that only valid values are stored.
	* `Get() T`: Returns the current value.

## Built-in Validators

| Validator | Type | Tag | Params (defaults) | Description |
| --- | --- | --- | --- | --- |
| **Generic** |  |  |  |  |
| `CmpRangeValidator[T]` | `cmp.Ordered` | - | `min`, `max` | Inclusive range for ordered types. |
| `NonZeroValidator[T]` | `any` | - | - | Value is not zero. |
| **Ints** |  |  |  |  |
| `IntRangeValidator` | `int` | `rangeint` | `min`, `max` | Inclusive int range. |
| `MinIntValidator` | `int` | `minint` | `min` | Int greater than or equal to `min`. |
| `MaxIntValidator` | `int` | `maxint` | `max` | Int less than or equal to `max`. |
| `NonNegativeIntValidator` | `int` | `posint` | - | Int is non-negative. |
| `NonPositiveIntValidator` | `int` | `negint` | - | Int is non-positive. |
| `NonZeroIntValidator` | `int` | `!zeroint` (alias: `nonzeroint`) | - | Int is not zero. |
| `OneOfIntValidator` | `int` | `oneofint` | `values` | Int is in `values` (pipe-separated list). |
| **Float64** |  |  |  |  |
| `Float64RangeValidator` | `float64` | `rangefloat` | `min`, `max` | Inclusive float64 range. |
| `MinFloat64Validator` | `float64` | `minfloat` | `min` | Float64 greater than or equal to `min`. |
| `MaxFloat64Validator` | `float64` | `maxfloat` | `max` | Float64 less than or equal to `max`. |
| `NonNegativeFloat64Validator` | `float64` | `posfloat` | - | Float64 is non-negative. |
| `NonPositiveFloat64Validator` | `float64` | `negfloat` | - | Float64 is non-positive. |
| `NonZeroFloat64Validator` | `float64` | `!zerofloat` (alias: `nonzerofloat`) | - | Float64 is not zero. |
| `OneOfFloat64Validator` | `float64` | `oneoffloat` | `values` | Float64 is in `values` (pipe-separated list). |
| **Strings** |  |  |  |  |
| `UrlValidator` | `string` | `url` | - | Valid URL (absolute). |
| `EmailValidator` | `string` | `email` | - | Valid email address. |
| `NonEmptyStringValidator` | `string` | `!empty` (alias: `nonempty`) | - | String is not empty. |
| `MinLengthValidator` | `string` | `min` | `size` | String length >= `size`. |
| `MaxLengthValidator` | `string` | `max` | `size` | String length <= `size`. |
| `LengthRangeValidator` | `string` | `len` | `min`, `max` | String length in inclusive range. |
| `RegexValidator` | `string` | `regex` | `pattern` | String matches regex. |
| `PrefixValidator` | `string` | `prefix` | `value` | String has prefix. |
| `SuffixValidator` | `string` | `suffix` | `value` | String has suffix. |
| `ContainsValidator` | `string` | `contains` | `value` | String contains substring. |
| `OneOfStringValidator` | `string` | `oneof` | `values` | String is in `values` (pipe-separated list). |
| `AlphaNumericValidator` | `string` | `alphanum` | - | String is alphanumeric. |
| `MACAddressValidator` | `string` | `mac` | - | Valid MAC address. |
| `IpValidator` | `string` | `ip` | - | Valid IP address. |
| `IPv4Validator` | `string` | `ipv4` | - | Valid IPv4 address. |
| `IPv6Validator` | `string` | `ipv6` | - | Valid IPv6 address. |
| `HostnameValidator` | `string` | `hostname` | - | Valid hostname. |
| `IPCIDRValidator` | `string` | `cidr` | - | Valid CIDR notation. |
| `UUIDValidator` | `string` | `uuid` | `version` (`4`) | RFC 4122 UUID with optional version. |
| `Base64Validator` | `string` | `base64` | - | Valid base64 (standard or raw). |
| `HexValidator` | `string` | `hex` | - | Valid hex string (optional `0x`). |
| `XMLValidator` | `string` | `xml` | - | Well-formed XML with at least one element. |
| `JSONValidator` | `string` | `json` | - | Valid JSON. |
| `TimeValidator` | `string` | `time` | `format` (`RFC3339`) | Valid time for layout (built-in names or raw layout). |
| **Time** |  |  |  |  |
| `NonZeroTimeValidator` | `time.Time` | `!zerotime` (alias: `nonzerotime`) | - | Time is not zero. |
| `TimeBeforeValidator` | `time.Time` | `beforetime` | `before` | Time is before the configured time (RFC3339). |
| `TimeAfterValidator` | `time.Time` | `aftertime` | `after` | Time is after the configured time (RFC3339). |
| `TimeBetweenValidator` | `time.Time` | `betweentime` | `start`, `end` | Time is within the inclusive range (RFC3339). |
| **Duration** |  |  |  |  |
| `PositiveDurationValidator` | `time.Duration` | `posduration` | - | Duration is positive. |
| `NonZeroDurationValidator` | `time.Duration` | `!zeroduration` (alias: `nonzeroduration`) | - | Duration is not zero. |
| **IP** |  |  |  |  |
| `NonZeroIPValidator` | `net.IP` | `!zeroip` (alias: `nonzeroip`) | - | IP is not zero or unspecified. |
| `IPRangeValidator` | `net.IP` | `iprange` | `start`, `end` | IP is within the inclusive range. |
| **URL** |  |  |  |  |
| `NonZeroURLValidator` | `url.URL` | `!zerourl` (alias: `nonzerourl`) | - | URL is not zero. |

## License

This project is licensed under the MIT License – see the [LICENSE](https://github.com/tedla-brandsema/valex/blob/main/LICENSE) file for details.
