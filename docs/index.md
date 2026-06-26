# Valex documentation

Valex is a validation library for Go. It pairs a small, dependency-light
validation *engine* with opt-in packages for a ready-made directive catalog and
HTTP form binding, so you depend only on what you use.

The library is split into three packages:

| Package | What it gives you |
| --- | --- |
| `valex` | the engine: `Validator[T]`, `ValidatorFunc[T]`, `ValidatedValue[T]`, `MustValidate`, the `val` struct tag (`ValidateStruct`, `RegisterDirective`), and re-exported error types. |
| `valex/validators` | a catalog of ready-made `val` directives (ranges, lengths, URLs, emails, IPs, time, JSON/XML, regex, …), registered opt-in. |
| `valex/forms` | binds `net/http` request values into structs and validates them, kept separate so the core never imports `net/http`. |

- [Quick start](quick-start.md) — install, register a directive, validate a struct.
- [Programmatic validation](programmatic.md) — `Validator[T]`, `ValidatorFunc[T]`, `ValidatedValue[T]`, and `MustValidate`.
- [Struct-tag validation](struct-tags.md) — the `val` tag, `ValidateStruct`, the validators catalog, and custom directives.
- [HTTP forms](forms.md) — bind and validate `net/http` requests with `valex/forms`.
- [Errors](errors.md) — the re-exported typed error model and how to inspect it with `errors.As`.

Runnable programs live in [examples/](../examples/). Go testable examples that
render on [pkg.go.dev](https://pkg.go.dev/github.com/tedla-brandsema/valex) live
in `example_test.go`.

Valex builds on [tagex](https://github.com/tedla-brandsema/tagex), the struct-tag
processor that does the reflection and directive dispatch underneath. You rarely
need to touch it directly — valex re-exports the error types — but the lifecycle
hooks ([Before/Success/Failure](forms.md#lifecycle-hooks)) are tagex's.
