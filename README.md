# Valex

[![Go Reference](https://pkg.go.dev/badge/github.com/tedla-brandsema/valex.svg)](https://pkg.go.dev/github.com/tedla-brandsema/valex)
[![CI](https://github.com/tedla-brandsema/valex/actions/workflows/ci.yml/badge.svg)](https://github.com/tedla-brandsema/valex/actions/workflows/ci.yml)

*Valex* is an extensible validation library for Go. It pairs a small, dependency-light
validation **engine** with opt-in packages for a ready-made directive catalog and HTTP
form binding, so you depend only on what you use.

## Packages

| Import | Responsibility |
| --- | --- |
| `github.com/tedla-brandsema/valex` | The engine: the `Validator[T]` interface and `ValidatorFunc[T]` adapter, the `ValidatedValue[T]` wrapper, `MustValidate`, the `val` struct tag (`ValidateStruct`), `RegisterDirective`, and re-exported error types. |
| `github.com/tedla-brandsema/valex/validators` | A catalog of ready-made `val` directives (ranges, lengths, URLs, emails, IPs, time, JSON/XML, regex, …). Directives are **opt-in** — you register the ones you want. |
| `github.com/tedla-brandsema/valex/forms` | Bind `net/http` request values into structs and validate them. Kept separate so the core engine never imports `net/http`. |

## Features

* **Generic validators** — define type-safe validators via the `Validator[T]` interface or the `ValidatorFunc[T]` adapter.
* **Validated value wrapper** — `ValidatedValue[T]` only stores values that pass validation.
* **Tag-based validation** — validate struct fields with the `val` tag and `ValidateStruct`.
* **Opt-in directive catalog** — register only the directives you need from `valex/validators`.
* **Custom directives** — extend the `val` tag with `RegisterDirective`.
* **HTTP form binding** — parse and validate requests with `valex/forms`.
* **Inspectable errors** — error types are re-exported from the engine, so you handle them without importing `tagex`.

## Installation

```
go get -u github.com/tedla-brandsema/valex@latest
```

## Programmatic validation

Implement the `Validator[T]` interface for your type, or adapt a function with
`ValidatorFunc[T]`, and use `ValidatedValue[T]` for guarded assignment:

```go
package main

import (
	"fmt"

	"github.com/tedla-brandsema/valex"
)

func main() {
	// A quick validator from a function.
	nonEmpty := valex.ValidatorFunc[string](func(val string) error {
		if val == "" {
			return fmt.Errorf("string cannot be empty")
		}
		return nil
	})

	vv := valex.ValidatedValue[string]{Validator: nonEmpty}
	if err := vv.Set("hello world"); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("validated value:", vv.Get())
}
```

## Tag-based validation

The engine ships **no** directives of its own. Register the ones you need from
`valex/validators` (typically in `init`), then call `ValidateStruct`:

```go
package main

import (
	"fmt"

	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/valex/validators"
)

func init() {
	valex.RegisterDirective(&validators.MinLengthValidator{})
	valex.RegisterDirective(&validators.EmailValidator{})
	valex.RegisterDirective(&validators.IntRangeValidator{})
}

type User struct {
	Name  string `val:"min,size=3"`
	Email string `val:"email"`
	Age   int    `val:"rangeint,min=0,max=120"`
}

func main() {
	if err := valex.ValidateStruct(&User{Name: "Al", Email: "invalid", Age: 200}); err != nil {
		fmt.Println(err)
	}
}
```

See the [`validators`](https://pkg.go.dev/github.com/tedla-brandsema/valex/validators)
package for the full catalog, or the [table below](#built-in-directives).

## Custom directives

A directive is any `tagex.Directive[T]` — implement `Name`, `Mode`, and `Handle`,
then register it:

```go
package main

import (
	"fmt"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
)

type EvenDirective struct{}

func (d *EvenDirective) Name() string              { return "even" }
func (d *EvenDirective) Mode() tagex.DirectiveMode { return tagex.EvalMode }
func (d *EvenDirective) Handle(val int) (int, error) {
	if val%2 != 0 {
		return val, fmt.Errorf("value %d is not even", val)
	}
	return val, nil
}

func main() {
	valex.RegisterDirective(&EvenDirective{})

	type Item struct {
		Count int `val:"even"`
	}
	if err := valex.ValidateStruct(&Item{Count: 3}); err != nil {
		fmt.Println(err)
	}
}
```

## HTTP form validation

`valex/forms` binds request values into a struct using `field` tags, then validates
the `val` tags. `forms.New` calls `request.ParseForm`, which reads both POST bodies
and URL query parameters (so GET requests work too). `forms.Validate` is a convenience
wrapper that returns a `*forms.Error` carrying an HTTP status code; `forms.Bind` binds
without validating when you are outside an HTTP handler.

```go
type Signup struct {
	Name  string `field:"name" val:"min,size=3"`
	Email string `field:"email" val:"email"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var in Signup
	if err := forms.Validate(r, &in); err != nil {
		var ferr *forms.Error
		errors.As(err, &ferr)
		http.Error(w, err.Error(), ferr.StatusCode())
		return
	}
	// ... use in
}
```

## Error handling

`ValidateStruct` and the `forms` helpers return errors you can inspect with
`errors.As` / `errors.Is`. The engine re-exports the underlying error types, so you
do **not** need to import `tagex`:

```go
err := valex.ValidateStruct(&User{ /* ... */ })
if err != nil {
	var conv *valex.ConversionError
	switch {
	case errors.As(err, &conv):
		// a parameter value could not be converted
	case errors.Is(err, valex.ErrNoValidator):
		// ...
	}
}
```

## Examples

Runnable programs in [examples/](examples/) — run one with `go run ./examples/<name>`:

- [programmatic](examples/programmatic/) — validate values in code with `ValidatorFunc` and `ValidatedValue`, no tags.
- [validate-struct](examples/validate-struct/) — register catalog directives and validate a struct with the `val` tag.
- [custom-directive](examples/custom-directive/) — extend the `val` tag with your own directive.
- [forms](examples/forms/) — bind and validate an `net/http` request with `valex/forms`.

## Documentation

Full documentation is in [docs/](docs/index.md):

- [Quick start](docs/quick-start.md) — install, register a directive, validate a struct.
- [Programmatic validation](docs/programmatic.md) — `Validator[T]`, `ValidatorFunc[T]`, `ValidatedValue[T]`, `MustValidate`.
- [Struct-tag validation](docs/struct-tags.md) — the `val` tag, the validators catalog, and custom directives.
- [HTTP forms](docs/forms.md) — bind and validate `net/http` requests.
- [Errors](docs/errors.md) — the re-exported typed error model.

Package reference and Go testable examples render on
[pkg.go.dev](https://pkg.go.dev/github.com/tedla-brandsema/valex). Working notes
and deferred decisions live in [TODO.md](TODO.md).

## Built-in directives

From `github.com/tedla-brandsema/valex/validators`. Register each with
`valex.RegisterDirective(&XxxValidator{})` before validating.

| Validator | Type | Tag | Params (defaults) | Description |
| --- | --- | --- | --- | --- |
| **Generic (programmatic, no tag)** |  |  |  |  |
| `CmpRangeValidator[T]` | `cmp.Ordered` | - | `Min`, `Max` | Inclusive range for ordered types. |
| `NonZeroValidator[T]` | `any` | - | - | Value is not the zero value. |
| `CompositeValidator[T]` | `cmp.Ordered` | - | `Validators` | Runs several validators in order. |
| **Ints** |  |  |  |  |
| `IntRangeValidator` | `int` | `rangeint` | `min`, `max` | Inclusive int range. |
| `MinIntValidator` | `int` | `minint` | `min` | Int `>= min`. |
| `MaxIntValidator` | `int` | `maxint` | `max` | Int `<= max`. |
| `NonNegativeIntValidator` | `int` | `posint` | - | Int is non-negative. |
| `NonPositiveIntValidator` | `int` | `negint` | - | Int is non-positive. |
| `NonZeroIntValidator` | `int` | `!zeroint` | - | Int is not zero. |
| `OneOfIntValidator` | `int` | `oneofint` | `values` | Int is in `values` (pipe-separated). |
| **Float64** |  |  |  |  |
| `Float64RangeValidator` | `float64` | `rangefloat` | `min`, `max` | Inclusive float64 range. |
| `MinFloat64Validator` | `float64` | `minfloat` | `min` | Float64 `>= min`. |
| `MaxFloat64Validator` | `float64` | `maxfloat` | `max` | Float64 `<= max`. |
| `NonNegativeFloat64Validator` | `float64` | `posfloat` | - | Float64 is non-negative. |
| `NonPositiveFloat64Validator` | `float64` | `negfloat` | - | Float64 is non-positive. |
| `NonZeroFloat64Validator` | `float64` | `!zerofloat` | - | Float64 is not zero. |
| `OneOfFloat64Validator` | `float64` | `oneoffloat` | `values` | Float64 is in `values` (pipe-separated). |
| **Strings** |  |  |  |  |
| `UrlValidator` | `string` | `url` | - | Valid absolute URL. |
| `EmailValidator` | `string` | `email` | - | Valid email address. |
| `NonEmptyStringValidator` | `string` | `!empty` | - | String is not empty. |
| `MinLengthValidator` | `string` | `min` | `size` | String length `>= size`. |
| `MaxLengthValidator` | `string` | `max` | `size` | String length `<= size`. |
| `LengthRangeValidator` | `string` | `len` | `min`, `max` | String length in inclusive range. |
| `RegexValidator` | `string` | `regex` | `pattern` | String matches regex. |
| `PrefixValidator` | `string` | `prefix` | `value` | String has prefix. |
| `SuffixValidator` | `string` | `suffix` | `value` | String has suffix. |
| `ContainsValidator` | `string` | `contains` | `value` | String contains substring. |
| `OneOfStringValidator` | `string` | `oneof` | `values` | String is in `values` (pipe-separated). |
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
| `TimeValidator` | `string` | `time` | `format` (`RFC3339`) | Valid time for layout (built-in name or raw layout). |
| **Time** |  |  |  |  |
| `NonZeroTimeValidator` | `time.Time` | `!zerotime` | - | Time is not zero. |
| `TimeBeforeValidator` | `time.Time` | `beforetime` | `before` | Time is before the configured time (RFC3339). |
| `TimeAfterValidator` | `time.Time` | `aftertime` | `after` | Time is after the configured time (RFC3339). |
| `TimeBetweenValidator` | `time.Time` | `betweentime` | `start`, `end` | Time is within the inclusive range (RFC3339). |
| **Duration** |  |  |  |  |
| `PositiveDurationValidator` | `time.Duration` | `posduration` | - | Duration is positive. |
| `NonZeroDurationValidator` | `time.Duration` | `!zeroduration` | - | Duration is not zero. |
| **IP** |  |  |  |  |
| `NonZeroIPValidator` | `net.IP` | `!zeroip` | - | IP is not zero or unspecified. |
| `IPRangeValidator` | `net.IP` | `iprange` | `start`, `end` | IP is within the inclusive range. |
| **URL** |  |  |  |  |
| `NonZeroURLValidator` | `url.URL` | `!zerourl` | - | URL is not the zero value. |

## License

This project is licensed under the MIT License – see the [LICENSE](https://github.com/tedla-brandsema/valex/blob/main/LICENSE) file for details.
