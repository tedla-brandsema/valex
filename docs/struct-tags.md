# Struct-tag validation

Tag a struct field with `val` and `ValidateStruct` walks the struct, runs the
named directive on each tagged field, and returns the first failure.

```go
type User struct {
	Name  string `val:"min,size=3"`
	Email string `val:"email"`
	Age   int    `val:"rangeint,min=0,max=120"`
}

err := valex.ValidateStruct(&User{Name: "Al", Email: "x", Age: 200})
```

`ValidateStruct` takes a **pointer to a struct**. It recurses into nested
structs, so a `val` tag deep in the tree is found and the error path points to
it. It returns `nil` on success or a typed [error](errors.md) on the first
failure.

## The tag grammar

```
val:"<directive>,<arg>=<value>,<arg>=<value>"
```

- The first token is the directive name (`min`, `email`, `rangeint`).
- The rest are `key=value` parameters consumed by that directive.

Where a directive takes a list, the values are pipe-separated:
`val:"oneof,values=draft|sent|paid"`.

## Registering directives

The engine ships **no** directives. Register the ones you use — once, at startup:

```go
func init() {
	valex.MustRegisterDirective(&validators.EmailValidator{})
	valex.MustRegisterDirective(&validators.IntRangeValidator{})
}
```

`MustRegisterDirective[T]` infers `T` from the directive, so the call site needs
no type argument. It panics if the directive's name is blank
(`*EmptyDirectiveNameError`) or already registered (`*DuplicateDirectiveError`) —
both setup-time programming mistakes, so failing fast at startup is what you
want. Use `RegisterDirective`, which returns that error instead of panicking, if
you register dynamically and need to handle it. Registering is the only operation
that mutates the `val` tag's registry; once registered, `ValidateStruct` is safe
to call from many goroutines (see [Concurrency](#concurrency)).

## The catalog

`valex/validators` provides ready-made directives, grouped by the field type they
validate. The **Registers** column is the type you pass to
`valex.RegisterDirective`.

### int

| Tag | Registers | Params | Checks |
| --- | --- | --- | --- |
| `rangeint` | `IntRangeValidator` | `min`, `max` | inclusive range |
| `minint` | `MinIntValidator` | `min` | `value >= min` |
| `maxint` | `MaxIntValidator` | `max` | `value <= max` |
| `posint` | `NonNegativeIntValidator` | — | non-negative |
| `negint` | `NonPositiveIntValidator` | — | non-positive |
| `!zeroint` | `NonZeroIntValidator` | — | not zero |
| `oneofint` | `OneOfIntValidator` | `values` | one of a pipe-separated list |

### float64

| Tag | Registers | Params | Checks |
| --- | --- | --- | --- |
| `rangefloat` | `Float64RangeValidator` | `min`, `max` | inclusive range |
| `minfloat` | `MinFloat64Validator` | `min` | `value >= min` |
| `maxfloat` | `MaxFloat64Validator` | `max` | `value <= max` |
| `posfloat` | `NonNegativeFloat64Validator` | — | non-negative |
| `negfloat` | `NonPositiveFloat64Validator` | — | non-positive |
| `!zerofloat` | `NonZeroFloat64Validator` | — | not zero |
| `oneoffloat` | `OneOfFloat64Validator` | `values` | one of a pipe-separated list |

### string

| Tag | Registers | Params | Checks |
| --- | --- | --- | --- |
| `url` | `UrlValidator` | — | valid absolute URL |
| `email` | `EmailValidator` | — | valid email address |
| `!empty` | `NonEmptyStringValidator` | — | non-empty |
| `min` | `MinLengthValidator` | `size` | `length >= size` |
| `max` | `MaxLengthValidator` | `size` | `length <= size` |
| `len` | `LengthRangeValidator` | `min`, `max` | length within `[min, max]` |
| `regex` | `RegexValidator` | `pattern` | matches regular expression |
| `prefix` | `PrefixValidator` | `value` | has prefix |
| `suffix` | `SuffixValidator` | `value` | has suffix |
| `contains` | `ContainsValidator` | `value` | contains substring |
| `oneof` | `OneOfStringValidator` | `values` | one of a pipe-separated list |
| `alphanum` | `AlphaNumericValidator` | — | alphanumeric |
| `mac` | `MACAddressValidator` | — | valid MAC address |
| `ip` | `IpValidator` | — | valid IP address |
| `ipv4` | `IPv4Validator` | — | valid IPv4 address |
| `ipv6` | `IPv6Validator` | — | valid IPv6 address |
| `hostname` | `HostnameValidator` | — | valid hostname |
| `cidr` | `IPCIDRValidator` | — | valid CIDR notation |
| `uuid` | `UUIDValidator` | `version` (`4`) | RFC 4122 UUID, optional version |
| `base64` | `Base64Validator` | — | valid base64 (standard or raw) |
| `hex` | `HexValidator` | — | valid hex (optional `0x`) |
| `xml` | `XMLValidator` | — | well-formed XML |
| `json` | `JSONValidator` | — | valid JSON |
| `time` | `TimeValidator` | `format` (`RFC3339`) | valid time for the layout |

### time.Time, time.Duration, net.IP, url.URL

| Tag | Registers | Params | Checks |
| --- | --- | --- | --- |
| `!zerotime` | `NonZeroTimeValidator` | — | time is not zero |
| `beforetime` | `TimeBeforeValidator` | `before` | before the given RFC3339 time |
| `aftertime` | `TimeAfterValidator` | `after` | after the given RFC3339 time |
| `betweentime` | `TimeBetweenValidator` | `start`, `end` | within `[start, end]` (RFC3339) |
| `posduration` | `PositiveDurationValidator` | — | duration is positive |
| `!zeroduration` | `NonZeroDurationValidator` | — | duration is not zero |
| `!zeroip` | `NonZeroIPValidator` | — | IP is not zero/unspecified |
| `iprange` | `IPRangeValidator` | `start`, `end` | IP within `[start, end]` |
| `!zerourl` | `NonZeroURLValidator` | — | URL is not the zero value |

## Custom directives

A directive is any `tagex.Directive[T]` — implement `Name`, `Mode`, and `Handle`
for the field type `T` you handle, then register it:

```go
type EvenDirective struct{}

func (*EvenDirective) Name() string              { return "even" }
func (*EvenDirective) Mode() tagex.DirectiveMode { return tagex.EvalMode }
func (*EvenDirective) Handle(n int) (int, error) {
	if n%2 != 0 {
		return n, fmt.Errorf("value %d is not even", n)
	}
	return n, nil
}

valex.MustRegisterDirective(&EvenDirective{})
// type Ticket struct { Seats int `val:"even"` }
```

- `Name()` is the tag name (`val:"even"`).
- `Mode()` is `tagex.EvalMode` to validate, or `tagex.MutMode` to write `Handle`'s
  return value back to the field (normalization).
- Parameters are declared by tagging the directive's own fields with `param`
  (a `Size int` field tagged `param:"size"`), which valex fills from the tag
  args before `Handle` runs.

For mutation, parameters, and conversion details, see the
[tagex directive](https://github.com/tedla-brandsema/tagex/blob/main/docs/directives.md)
and [parameter](https://github.com/tedla-brandsema/tagex/blob/main/docs/parameters.md)
guides — valex registers and runs `tagex.Directive` values unchanged.

## Multiple tags in one pass

`ValidateStruct` accepts extra `*tagex.Tag` values to process alongside `val` in
a single walk:

```go
err := valex.ValidateStruct(&data, otherTag)
```

## Concurrency

`RegisterDirective` and `ValidateStruct` are safe for concurrent use. Register
directives once at startup (typically in `init`), then validate from any number
of goroutines. Registering while other goroutines validate is safe but unusual.

## One global registry

The `val` directives live in a single, process-wide registry — the one behind
`RegisterDirective`, `MustRegisterDirective`, and `ValidateStruct`. Think of it
the way you think of `flag.CommandLine` or `http.DefaultServeMux`: **it belongs to
the application**, not to any one library.

**If you are writing an application**, this is all you need. Register your
directives once at startup with `MustRegisterDirective` and validate anywhere.
Registration fails loudly on a duplicate name (`*DuplicateDirectiveError`), so a
double-registration is caught at boot instead of silently shadowing.

**If you are writing a library** that validates internally, **do not register on
the global registry.** It is shared with the importing program and every other
library in the binary, so if two of them register the same directive name — say
two libraries that both `MustRegisterDirective(&validators.EmailValidator{})` —
the second one panics at startup. Bring your own registry instead: the catalog
validators are plain `tagex.Directive` values, so they register on any tag.

```go
reg := tagex.NewTag("val")
tagex.MustRegisterDirective(reg, &validators.EmailValidator{})

err := tagex.ProcessStruct(&user, reg) // isolated; the global registry is untouched
```

Validate with `tagex.ProcessStruct(&data, reg)`, **not** `valex.ValidateStruct(&data, reg)`:
`ValidateStruct` always *also* applies the global `val` registry, and since both
tags share the key `val`, every field would be processed by both — and the global
pass would fail on any directive only `reg` knows. The extra-tags parameter of
`ValidateStruct` (see [Multiple tags in one pass](#multiple-tags-in-one-pass)) is
for tags with a *different* key, not for same-key isolation.
