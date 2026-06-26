# Programmatic validation

The engine's core is one interface — no tags, no reflection:

```go
type Validator[T any] interface {
	Validate(val T) error
}
```

A `nil` error means the value is valid. Everything else in this page is built on
that one method.

## ValidatorFunc

`ValidatorFunc[T]` adapts a plain function to `Validator[T]`, so you don't need a
named type for a one-off rule:

```go
nonEmpty := valex.ValidatorFunc[string](func(s string) error {
	if s == "" {
		return errors.New("must not be empty")
	}
	return nil
})

err := nonEmpty.Validate("") // "must not be empty"
```

## ValidatedValue

`ValidatedValue[T]` wraps a value behind its validator: `Set` stores the value
only if it passes, and leaves the previous value untouched on failure.

```go
positive := valex.ValidatorFunc[int](func(n int) error {
	if n <= 0 {
		return errors.New("must be positive")
	}
	return nil
})

v := valex.ValidatedValue[int]{Validator: positive}
v.Set(42)        // nil; stored
v.Set(-1)        // "must be positive"; not stored
fmt.Println(v.Get()) // 42
```

`Set` returns `ErrNoValidator` if no `Validator` is configured. `Get` returns the
stored value and `String` formats it with `%v`.

## MustValidate

`MustValidate` validates and returns the value, or **panics** if it fails. Use it
for values you treat as invariants — config parsed at startup, test fixtures —
where an invalid value is a programming error, not a runtime condition:

```go
port := valex.MustValidate(8080, inRange) // returns 8080, or panics
```

## Ready-made validators

The `valex/validators` package also exposes a few `Validator[T]` values you can
use directly, without the `val` tag:

- `CmpRangeValidator[T cmp.Ordered]` — inclusive `[Min, Max]` range for any ordered type.
- `NonZeroValidator[T]` — the value is not its zero value.
- `CompositeValidator[T cmp.Ordered]` — runs several validators in order, returning the first failure.

```go
ageOK := validators.CmpRangeValidator[int]{Min: 0, Max: 120}
err := ageOK.Validate(200) // out of range
```

These complement the tag directives: the same `validators` package registers
directives for the `val` tag ([struct-tags.md](struct-tags.md)) and offers these
for direct use.
