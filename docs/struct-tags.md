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

## Chaining directives

Apply several directives to one field by separating them with `;`. They run
**left to right**, and each `MutMode` directive's written-back value is what the
next one sees:

```go
type Login struct {
	User string `val:"alphanum;min,size=3"`
}
```

`alphanum` and `min` are both checks, so order here only decides which failure
surfaces first. Order becomes significant once a `MutMode` (normalizing)
directive is in the chain: a custom `trim;min,size=3` validates the *trimmed*
length, while `min,size=3;trim` validates the raw one. The
[chained example](../examples/chained/) is a runnable `trim;lower;max` pipeline.

A chain **stops at the first failing segment** and reports that one error; later
segments do not run. Each segment gets its own per-call copy of the directive, so
chained directives never share parameter state. Stray separators are ignored, so
a leading, doubled, or trailing `;` (`;min,size=3;;`) is harmless.

Two things to know:

- **Reserved characters.** `,` separates parameters and `;` chains directives. To
  use either literally in a parameter value — or `=`, or significant
  leading/trailing whitespace — quote the value (see below).
- **Partial mutation under `ValidateStructAll`.** Because a chain stops mid-way on
  failure, a `MutMode` segment that already ran has still written to the field.
  `ValidateStructAll` records the error and continues to the other fields, leaving
  that field partially transformed. If you need all-or-nothing field mutation,
  validate before you mutate (`min,size=3;trim`, not `trim;min,size=3`).

## Quoting parameter values

A parameter value is delimited by the reserved characters around it: only the
first `=` in a pair splits key from value (so `pattern=a=b` gives value `a=b`),
but a `,` ends the parameter and a `;` ends the directive. To put `,`, `;`, `=`,
or significant leading/trailing whitespace **inside** a value, wrap it in single
quotes:

```go
type Rule struct {
	Code string `val:"regex,pattern='[a-z]{1,3}'"` // value is [a-z]{1,3}, comma and all
}
```

Single quotes are used because the tag value is already delimited by double
quotes (`val:"..."`). Inside the quotes `,`, `;`, and `=` are literal; double an
interior quote (`''`) for a literal `'`, and a quoted empty value (`sep=''`) is an
explicit empty string.

A literal backslash must be **doubled** in the struct tag itself
(`pattern='^\\d{1,3}$'`) — Go unquotes the tag value before valex sees it, so a
lone `\d` is rejected by `go vet` as a bad tag. This is Go's struct-tag rule, not
valex's.

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

## Registries

The package-level `RegisterDirective`, `MustRegisterDirective`, and
`ValidateStruct` share one process-wide **default registry**. Think of it the way
you think of `flag.CommandLine` or `http.DefaultServeMux`: **it belongs to the
application**, not to any one library.

**If you are writing an application**, use the package-level functions: register
your directives once at startup with `MustRegisterDirective` and validate
anywhere. Registration fails loudly on a duplicate name
(`*DuplicateDirectiveError`), so a double-registration is caught at boot instead
of silently shadowing.

**If you are writing a library, or you need test isolation**, create your own
registry with `NewRegistry` instead of touching the global. Each registry has an
independent directive set, so two of them can hold the same directive name
without colliding — and a test can spin up a fresh one per case:

```go
reg := valex.NewRegistry()
valex.MustRegisterDirectiveTo(reg, &validators.EmailValidator{})

err := reg.ValidateStruct(&user) // uses only reg's directives
```

Registration on a registry is the free function `RegisterDirectiveTo` /
`MustRegisterDirectiveTo` (rather than a method) because Go methods can't have
type parameters. For isolated **form** validation, pass the registry to `forms`
— one call, still returning a `*forms.Error` with an HTTP status:

```go
var in Signup
err := forms.ValidateWith(r, &in, reg) // bind + validate against reg
```
